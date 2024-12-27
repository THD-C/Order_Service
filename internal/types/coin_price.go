package types

import (
	"github.com/shopspring/decimal"
	"sync"
)

type CoinPrice struct {
	CoinId string
	Prices map[string]decimal.Decimal
	mu     sync.Mutex
}

func NewCoinPrice(coinId string) *CoinPrice {
	return &CoinPrice{
		CoinId: coinId,
		Prices: make(map[string]decimal.Decimal),
	}
}

func (cp *CoinPrice) SetPrices(prices map[string]decimal.Decimal) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.Prices = prices
}
