package server

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	proto "order_service/generated/order"
	"order_service/internal/service/order_service"
)

const (
	UnsupportedOrderSide = "unsupported order side"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
	buyOrderService  order_service.OrderService
	sellOrderService order_service.OrderService
}

func NewOrderServer() *OrderServer {
	return &OrderServer{
		buyOrderService:  &order_service.BuyOrderService{},
		sellOrderService: &order_service.SellOrderService{},
	}
}

func (s *OrderServer) CreateOrder(
	ctx context.Context,
	req *proto.OrderDetails,
) (*proto.OrderDetails, error) {
	log.Info().Interface("request", req).Msg("Creating order")
	switch req.Side {
	case proto.OrderSide_ORDER_SIDE_BUY:
		return s.buyOrderService.CreateOrder(ctx, req)
	case proto.OrderSide_ORDER_SIDE_SELL:
		return s.sellOrderService.CreateOrder(ctx, req)
	default:
		log.Error().Interface("request", req).Msg(UnsupportedOrderSide)
		return nil, fmt.Errorf(UnsupportedOrderSide)
	}
}
