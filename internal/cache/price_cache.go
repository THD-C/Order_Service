package cache

import (
	"context"
	"order_service/internal/types"
	"sync"
	"time"
)

type PriceCache struct {
	prices sync.Map
}

func NewPriceCache() *PriceCache {
	return &PriceCache{}
}

func (pc *PriceCache) GetPrice(symbol string) (types.Price, bool) {
	value, exists := pc.prices.Load(symbol)
	if !exists {
		return types.Price{}, false
	}
	price, ok := value.(types.Price)
	if !ok {
		return types.Price{}, false
	}
	return price, true
}

func (pc *PriceCache) UpdatePrices(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			prices, err := fetchPrices()
			if err != nil {
				continue
			}
			for symbol, price := range prices {
				pc.prices.Store(symbol, price)
			}
		case <-ctx.Done():
			return
		}
	}
}

func fetchPrices() (map[string]types.Price, error) {
	// TODO: Tutaj będzie łączenie się do Arka i zbieranie cen
	return nil, nil
}
