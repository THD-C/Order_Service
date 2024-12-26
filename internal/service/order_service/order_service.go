package order_service

import (
	"context"
	proto "order_service/generated/order"
	"order_service/internal/types"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req *proto.OrderDetails) (*proto.OrderDetails, error)
	ExecuteOrder(ctx context.Context, req *types.Order) error
}
