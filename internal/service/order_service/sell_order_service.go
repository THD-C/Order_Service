package order_service

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	proto "order_service/generated/order"
	"order_service/internal/bussiness_errors"
	"order_service/internal/cache"
	"order_service/internal/types"
)

type SellOrderService struct{}

func (s *SellOrderService) processOrder(
	order *types.Order,
) error {
	log.Info().Interface("request", order).Msg("Processing sell order")

	fiatWallet, cryptoWallet, err := FetchAndValidateWallets(order)
	if err != nil {
		return err
	}

	btcAmount := order.Nominal.Div(order.Price)

	if cryptoWallet.Value.LessThan(btcAmount) {
		log.Error().Interface("request", order).Msg(bussiness_errors.MsgInsufficientCryptoCurrency)
		return bussiness_errors.NewCustomError(
			bussiness_errors.ErrInsufficientCryptoCurrency,
			bussiness_errors.MsgInsufficientCryptoCurrency,
		)
	}

	originalCryptoValue := cryptoWallet.Value

	cryptoWallet.Value = cryptoWallet.Value.Sub(btcAmount)
	if err = cache.SaveWallet(cryptoWallet); err != nil {
		log.Error().Err(err).Interface("request", order).Msg("Failed to update crypto wallet")
		return fmt.Errorf("failed to update crypto wallet: %v", err)
	}

	fiatWallet.Value = fiatWallet.Value.Add(order.Nominal.Mul(order.Price))
	if err = cache.SaveWallet(fiatWallet); err != nil {
		if rollbackErr := RollbackWallet(cryptoWallet, originalCryptoValue); rollbackErr != nil {
			return fmt.Errorf(
				"failed to update fiat wallet: %v, and rollback failed: %v",
				err,
				rollbackErr,
			)
		}
		return fmt.Errorf("failed to update fiat wallet: %v", err)
	}

	log.Info().Interface("request", order).Msg("Sell order processed successfully")
	order.Status = proto.OrderStatus_ORDER_STATUS_COMPLETED
	return nil
}

func (s *SellOrderService) CreateOrder(
	_ context.Context,
	req *proto.OrderDetails,
) (*proto.OrderDetails, error) {
	order := &types.Order{}
	if err := order.FromProto(req); err != nil {
		return nil, err
	}

	if err := s.processOrder(order); err != nil {
		return nil, err
	}

	return order.ToProto(), nil
}

func (s *SellOrderService) ExecuteOrder(
	_ context.Context,
	order *types.Order,
) error {
	return s.processOrder(order)
}
