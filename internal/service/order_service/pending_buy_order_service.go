package order_service

import (
	"context"
	proto "order_service/generated/order"
	"order_service/internal/bussiness_errors"
	"order_service/internal/cache"
	"order_service/internal/logger"
	"order_service/internal/types"
	"time"
)

type PendingBuyOrderService struct {
	orderCache       *cache.OrderCache
	priceCache       *cache.PriceCache
	buyOrderService  OrderService
	sellOrderService OrderService
}

func NewPendingBuyOrderService(
	buyOrderService OrderService,
	sellOrderService OrderService,
) *PendingBuyOrderService {
	return &PendingBuyOrderService{
		orderCache:       cache.GetOrderCache(),
		priceCache:       cache.GetPriceCache(),
		buyOrderService:  buyOrderService,
		sellOrderService: sellOrderService,
	}
}

func (s *PendingBuyOrderService) CreateOrder(
	_ context.Context,
	order *proto.OrderDetails,
) (*proto.OrderDetails, error) {
	var myOrder types.Order
	if err := myOrder.FromProto(order); err != nil {
		order.Status = proto.OrderStatus_ORDER_STATUS_REJECTED
		return order, err
	}

	_, err := s.orderCache.Get(myOrder.ID)
	if err == nil {
		return order, bussiness_errors.NewCustomError(
			bussiness_errors.ErrOrderAlreadyExists,
			bussiness_errors.MsgOrderAlreadyExists,
		)
	}

	fiatWallet, _ := cache.FetchWallet(myOrder.FiatWalletID)
	cryptoWallet, _ := cache.FetchWallet(myOrder.CryptoWalletID)

	pendingOrder := &types.PendingOrder{
		Order:          &myOrder,
		FiatCurrency:   fiatWallet.Currency,
		CryptoCurrency: cryptoWallet.Currency,
	}
	err = s.orderCache.Add(pendingOrder)
	if err != nil {
		order.Status = proto.OrderStatus_ORDER_STATUS_REJECTED
		return order, err
	}

	order.Status = proto.OrderStatus_ORDER_STATUS_ACCEPTED
	return order, nil
}

func (s *PendingBuyOrderService) UpdateOrder(
	_ context.Context,
	order *types.Order,
) error {
	cachedOrder, err := s.orderCache.Get(order.ID)
	if err != nil {
		return bussiness_errors.NewCustomError(
			bussiness_errors.ErrOrderNotFound,
			bussiness_errors.MsgOrderNotFound,
		)
	}

	cachedOrder.Order = order
	err = s.orderCache.Update(cachedOrder)
	if err != nil {
		return err
	}

	return nil
}

func (s *PendingBuyOrderService) DeleteOrder(
	_ context.Context,
	orderID string,
) error {
	_, err := s.orderCache.Get(orderID)
	if err != nil {
		return bussiness_errors.NewCustomError(
			bussiness_errors.ErrOrderNotFound,
			bussiness_errors.MsgOrderNotFound,
		)
	}

	err = s.orderCache.Delete(orderID)
	if err != nil {
		return err
	}

	return nil
}

func (s *PendingBuyOrderService) CheckAndExecuteOrders(
	ctx context.Context,
	interval time.Duration,
) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.executePendingOrders(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (s *PendingBuyOrderService) executePendingOrders(ctx context.Context) {
	s.executePendingBuyOrders(ctx)
	s.executePendingSellOrders(ctx)
}

func (s *PendingBuyOrderService) executePendingBuyOrders(ctx context.Context) {
	var err error
	log := logger.GetLogger()
	orders, err := s.orderCache.GetAll()
	if err != nil {
		return
	}

	for _, pendingOrder := range orders {
		if pendingOrder.Order.Side != proto.OrderSide_ORDER_SIDE_BUY {
			continue
		}

		priceCoin, exists := s.priceCache.GetPrice(pendingOrder.CryptoCurrency)
		if !exists {
			log.Error().Msgf("Could not find prices for %s", pendingOrder.CryptoCurrency)
			continue
		}

		price := priceCoin.Prices[pendingOrder.FiatCurrency]

		if price.LessThanOrEqual(pendingOrder.Order.Price) {
			pendingOrder.Order.Price = price
			err = s.buyOrderService.ExecuteOrder(ctx, pendingOrder.Order)
			if err != nil {
				log.Error().Err(err).Msg("Failed to execute order")
				continue
			}

			err = s.orderCache.Delete(pendingOrder.Order.ID)
			if err != nil {
				log.Error().Err(err).Msg("Failed to delete order from cache")
			}
		}
	}
}

func (s *PendingBuyOrderService) executePendingSellOrders(ctx context.Context) {
	var err error
	log := logger.GetLogger()
	orders, err := s.orderCache.GetAll()
	if err != nil {
		return
	}

	for _, pendingOrder := range orders {
		if pendingOrder.Order.Side != proto.OrderSide_ORDER_SIDE_SELL {
			continue
		}

		priceCoin, exists := s.priceCache.GetPrice(pendingOrder.CryptoCurrency)
		if !exists {
			log.Error().Msgf("Could not find prices for %s", pendingOrder.CryptoCurrency)
			continue
		}

		price := priceCoin.Prices[pendingOrder.FiatCurrency]

		if price.GreaterThanOrEqual(pendingOrder.Order.Price) {
			pendingOrder.Order.Price = price
			err = s.sellOrderService.ExecuteOrder(ctx, pendingOrder.Order)
			if err != nil {
				log.Error().Err(err).Msg("Failed to execute order")
				continue
			}

			err = s.orderCache.Delete(pendingOrder.Order.ID)
			if err != nil {
				log.Error().Err(err).Msg("Failed to delete order from cache")
			}
		}
	}
}
