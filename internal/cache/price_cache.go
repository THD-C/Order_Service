package cache

import (
	"context"
	"github.com/shopspring/decimal"
	"order_service/internal/client"
	"order_service/internal/logger"
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

func (pc *PriceCache) GetPrice(symbol string) (*types.CoinPrice, bool) {
	value, exists := pc.prices.Load(symbol)
	if !exists {
		return &types.CoinPrice{}, false
	}
	prices, ok := value.(map[string]decimal.Decimal)
	if !ok {
		return &types.CoinPrice{}, false
	}
	priceCoin := types.CoinPrice{
		CoinId: symbol,
		Prices: prices,
	}
	return &priceCoin, true
}

func (pc *PriceCache) UpdatePrices(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	log := logger.GetLogger()

	for {
		select {
		case <-ticker.C:
			prices, err := fetchPrices()
			if err != nil {
				log.Error().Err(err).Msg("Failed to fetch prices")
				continue
			}
			for _, price := range prices {
				pc.prices.Store(price.CoinId, price.Prices)
			}
		case <-ctx.Done():
			return
		}
	}
}

func fetchPrices() ([]*types.CoinPrice, error) {
	coinGeckoClient, err := client.GetCoinGeckoClient()

	if err != nil {
		return nil, err
	}

	prices, err := coinGeckoClient.GetAllCoinsPrice()
	if err != nil {
		return nil, err
	}

	return prices, nil
}
