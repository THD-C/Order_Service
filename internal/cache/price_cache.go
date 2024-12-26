package cache

import (
	"context"
	"github.com/shopspring/decimal"
	"sync"
	"time"
)

type PriceCache struct {
	prices sync.Map
}

func NewPriceCache() *PriceCache {
	return &PriceCache{}
}

func (pc *PriceCache) GetPrice(symbol string) (decimal.Decimal, bool) {
	value, exists := pc.prices.Load(symbol)
	if !exists {
		return decimal.Decimal{}, false
	}
	price, ok := value.(decimal.Decimal)
	if !ok {
		return decimal.Decimal{}, false
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

func fetchPrices() (map[string]decimal.Decimal, error) {
	// TODO: Tutaj będzie łączenie się do Arka i zbieranie cen
	return nil, nil
}
