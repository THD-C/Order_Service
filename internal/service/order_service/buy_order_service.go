package order_service

import (
	"context"
	"fmt"
	proto "order_service/generated/order"
	"order_service/internal/bussiness_errors"
	"order_service/internal/cache"
	"order_service/internal/client"
	"order_service/internal/logger"
	"order_service/internal/types"
)

type BuyOrderService struct{}

func NewBuyOrderService() *BuyOrderService {
	return &BuyOrderService{}
}

func (s *BuyOrderService) processOrder(
	order *types.Order,
) error {
	log := logger.GetLogger()
	log.Info().Interface("request", order).Msg("Processing buy order")

	fiatWallet, cryptoWallet, err := FetchAndValidateWallets(order)
	if err != nil {
		order.Status = proto.OrderStatus_ORDER_STATUS_REJECTED
		return err
	}

	totalCost := order.Nominal.Mul(order.Price)

	fiatWallet.Mutex.Lock()
	if fiatWallet.Value.LessThan(totalCost) {
		fiatWallet.Mutex.Unlock()
		log.Error().Interface("request", order).Msg(bussiness_errors.MsgInsufficientFiatCurrency)
		order.Status = proto.OrderStatus_ORDER_STATUS_REJECTED
		return bussiness_errors.NewCustomError(
			bussiness_errors.ErrInsufficientFiatCurrency,
			bussiness_errors.MsgInsufficientFiatCurrency,
		)
	}

	originalFiatValue := fiatWallet.Value

	fiatWallet.Value = fiatWallet.Value.Sub(totalCost)
	if err = cache.SaveWallet(fiatWallet); err != nil {
		log.Error().Err(err).Interface("request", order).Msg("Failed to update fiat wallet")
		order.Status = proto.OrderStatus_ORDER_STATUS_REJECTED
		return fmt.Errorf("failed to update fiat wallet: %v", err)
	}
	fiatWallet.Mutex.Unlock()

	btcAmount := order.Nominal

	cryptoWallet.Mutex.Lock()
	defer cryptoWallet.Mutex.Unlock()

	cryptoWallet.Value = cryptoWallet.Value.Add(btcAmount)
	if err = cache.SaveWallet(cryptoWallet); err != nil {
		order.Status = proto.OrderStatus_ORDER_STATUS_REJECTED
		if rollbackErr := RollbackWallet(fiatWallet, originalFiatValue); rollbackErr != nil {
			return fmt.Errorf(
				"failed to update crypto wallet: %v, and rollback failed: %v",
				err,
				rollbackErr,
			)
		}
		return fmt.Errorf("failed to update crypto wallet: %v", err)
	}

	dbManagerClient, _ := client.GetDBManagerClient()
	_ = dbManagerClient.UpdateWallet(fiatWallet)
	_ = dbManagerClient.UpdateWallet(cryptoWallet)

	log.Info().Interface("request", order).Msg("Buy order processed successfully")

	order.Status = proto.OrderStatus_ORDER_STATUS_COMPLETED
	return nil
}

func (s *BuyOrderService) CreateOrder(
	_ context.Context,
	req *proto.OrderDetails,
) (*proto.OrderDetails, error) {
	order := &types.Order{}
	if err := order.FromProto(req); err != nil {
		return nil, err
	}
	dbManagerClient, _ := client.GetDBManagerClient()

	err := s.processOrder(order)
	_ = dbManagerClient.UpdateOrder(order)

	if err != nil {
		return nil, err
	}

	return order.ToProto(), nil
}

func (s *BuyOrderService) ExecuteOrder(
	_ context.Context,
	order *types.Order,
) error {
	return s.processOrder(order)
}
