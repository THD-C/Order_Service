package cache

import (
	"context"
	"github.com/shopspring/decimal"
	"order_service/internal/types"
	"sync"
	"time"
)

var (
	priceCacheInstance *PriceCache
	priceCacheOnce     sync.Once
)

type PriceCache struct {
	prices sync.Map
}

func NewPriceCache() *PriceCache {
	priceCacheOnce.Do(
		func() {
			priceCacheInstance = &PriceCache{}
		},
	)
	return priceCacheInstance
}

func GetPriceCache() *PriceCache {
	return priceCacheInstance
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

func fetchPrices() ([]*types.CoinPrice, error) {
	// coinGeckoClient, err := client.GetCoinGeckoClient()
	//
	// if err != nil {
	// 	return nil, err
	// }
	return nil, nil
}
