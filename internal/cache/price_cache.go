package cache

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"net"
	"sync"
	"time"
)

type PriceCache struct {
	prices sync.Map
	conn   net.Conn
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
	// ctx, cancel := context.WithTimeout(context.Background(), config.GetConfig().CoingeckoServiceTimeout)
	// defer cancel()
	//
	// conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	// if err != nil {
	// 	return nil, err
	// }
	// defer conn.Close()
	//
	// if err != nil {
	// 	return nil, err
	// }
	//
	// result := make(map[string]decimal.Decimal)
	// for _, coin := range response.Coins {
	// 	price, err := decimal.NewFromString(coin.Price)
	// 	if err != nil {
	// 		return nil, errors.New("invalid price format for symbol: " + coin.Symbol)
	// 	}
	// 	result[coin.Symbol] = price
	// }

	// TODO: Trzeba poczekać na dodanie metody od Arka

	return nil, fmt.Errorf("tests")
}
