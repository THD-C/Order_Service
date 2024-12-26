package cache

import (
	"fmt"
	"order_service/internal/types"
	"sync"
)

type OrderCache struct {
	orderMap sync.Map
}

func NewOrderCache() *OrderCache {
	return &OrderCache{}
}

func (oc *OrderCache) Add(pendingOrder *types.PendingOrder) error {
	oc.orderMap.Store(pendingOrder.Order.ID, pendingOrder)
	return nil
}

func (oc *OrderCache) Get(orderID string) (*types.PendingOrder, error) {
	value, ok := oc.orderMap.Load(orderID)
	if !ok {
		return nil, fmt.Errorf("order not found")
	}

	pendingOrder, ok := value.(*types.PendingOrder)
	if !ok {
		return nil, fmt.Errorf("failed to assert order data")
	}

	return pendingOrder, nil
}

func (oc *OrderCache) GetAll() ([]*types.PendingOrder, error) {
	var orders []*types.PendingOrder
	oc.orderMap.Range(
		func(key, value interface{}) bool {
			pendingOrder, ok := value.(*types.PendingOrder)
			if !ok {
				return false
			}
			orders = append(orders, pendingOrder)
			return true
		},
	)
	return orders, nil
}

func (oc *OrderCache) Update(pendingOrder *types.PendingOrder) error {
	_, ok := oc.orderMap.Load(pendingOrder.Order.ID)
	if !ok {
		return fmt.Errorf("order not found")
	}

	oc.orderMap.Store(pendingOrder.Order.ID, pendingOrder)
	return nil
}

func (oc *OrderCache) Delete(orderID string) error {
	_, ok := oc.orderMap.Load(orderID)
	if !ok {
		return fmt.Errorf("order not found")
	}

	oc.orderMap.Delete(orderID)
	return nil
}
