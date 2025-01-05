package server

import (
	"context"
	"fmt"
	proto "order_service/generated/order"
	"order_service/internal/service/order_service"
	"order_service/internal/types"
	"time"
)

const (
	UnsupportedOrderSide = "unsupported order side"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
	buyOrderService        order_service.OrderService
	sellOrderService       order_service.OrderService
	pendingBuyOrderService *order_service.PendingBuyOrderService
}

func NewOrderServer() *OrderServer {
	buyService := &order_service.BuyOrderService{}
	sellService := &order_service.SellOrderService{}
	pendingBuyService := order_service.NewPendingBuyOrderService(buyService, sellService)
	go pendingBuyService.CheckAndExecuteOrders(context.Background(), time.Second*60)

	return &OrderServer{
		buyOrderService:        buyService,
		sellOrderService:       sellService,
		pendingBuyOrderService: pendingBuyService,
	}
}

func (s *OrderServer) CreateOrder(
	ctx context.Context,
	req *proto.OrderDetails,
) (*proto.OrderDetails, error) {
	if req.Type == proto.OrderType_ORDER_TYPE_PENDING {
		return s.pendingBuyOrderService.CreateOrder(ctx, req)
	}

	if req.Side == proto.OrderSide_ORDER_SIDE_BUY {
		return s.buyOrderService.CreateOrder(ctx, req)
	}

	if req.Side == proto.OrderSide_ORDER_SIDE_SELL {
		return s.sellOrderService.CreateOrder(ctx, req)
	}

	return nil, fmt.Errorf(UnsupportedOrderSide)
}

func (s *OrderServer) UpdateOrder(
	ctx context.Context,
	req *proto.OrderDetails,
) (*proto.OrderDetails, error) {
	var myOrder types.Order
	if err := myOrder.FromProto(req); err != nil {
		return nil, err
	}

	if req.Type == proto.OrderType_ORDER_TYPE_PENDING {
		err := s.pendingBuyOrderService.UpdateOrder(ctx, &myOrder)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf(UnsupportedOrderSide)
	}

	return req, nil
}

func (s *OrderServer) DeleteOrder(
	ctx context.Context,
	req *proto.OrderID,
) (*proto.OrderDetails, error) {
	err := s.pendingBuyOrderService.DeleteOrder(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &proto.OrderDetails{}, nil
}
