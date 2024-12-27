package order_service

import (
	"context"
	"github.com/rs/zerolog/log"
	proto "order_service/generated/order"
	"order_service/internal/bussiness_errors"
	"order_service/internal/cache"
	"order_service/internal/types"
	"time"
)

type PendingBuyOrderService struct {
	orderCache   *cache.OrderCache
	priceCache   *cache.PriceCache
	orderService OrderService
}

func NewPendingBuyOrderService(
	priceCache *cache.PriceCache,
	orderService OrderService,
) *PendingBuyOrderService {
	return &PendingBuyOrderService{
		orderCache:   cache.NewOrderCache(),
		priceCache:   priceCache,
		orderService: orderService,
	}
}

func (s *PendingBuyOrderService) CreateOrder(
	_ context.Context,
	order *types.Order,
	fiatCurrency, cryptoCurrency string,
) error {
	_, err := s.orderCache.Get(order.ID)
	if err == nil {
		return bussiness_errors.NewCustomError(
			bussiness_errors.ErrOrderAlreadyExists,
			bussiness_errors.MsgOrderAlreadyExists,
		)
	}

	// TODO: Check if price changed meantime

	pendingOrder := &types.PendingOrder{
		Order:          order,
		FiatCurrency:   fiatCurrency,
		CryptoCurrency: cryptoCurrency,
	}
	err = s.orderCache.Add(pendingOrder)
	if err != nil {
		return err
	}

	return nil
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
			// s.executePendingOrders(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (s *PendingBuyOrderService) executePendingOrders(ctx context.Context) {
	var err error
	orders, err := s.orderCache.GetAll()
	if err != nil {
		return
	}

	for _, pendingOrder := range orders {
		if pendingOrder.Order.Side != proto.OrderSide_ORDER_SIDE_BUY {
			continue
		}

		price, exists := s.priceCache.GetPrice(pendingOrder.CryptoCurrency)
		if !exists {
			continue
		}

		if price.LessThanOrEqual(pendingOrder.Order.Price) {
			err = s.orderService.ExecuteOrder(ctx, pendingOrder.Order)
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
