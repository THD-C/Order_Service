package types

import (
	"fmt"
	"github.com/shopspring/decimal"
	// proto "order_service/generated/coins"
)

type Coin struct {
	ID   string          `json:"id"`
	Buy  decimal.Decimal `json:"buy"`
	Sell decimal.Decimal `json:"sell"`
}

func NewCoin(id, buy, sell string) *Coin {
	buyVal, _ := decimal.NewFromString(buy)
	sellVal, _ := decimal.NewFromString(sell)

	return &Coin{
		ID:   id,
		Buy:  buyVal,
		Sell: sellVal,
	}
}

func (c *Coin) String() string {
	return fmt.Sprintf(
		"[Coin ID: %s, Buy: %s, Sell: %s]",
		c.ID,
		c.Buy,
		c.Sell,
	)
}

// func (c *Coin) FromProto(p *proto.Coin) error {
// 	buy, err := decimal.NewFromString(p.Buy)
// 	if err != nil || buy.LessThan(decimal.Zero) {
// 		return fmt.Errorf("wrong buy value")
// 	}
//
// 	sell, err := decimal.NewFromString(p.Sell)
// 	if err != nil || sell.LessThan(decimal.Zero) {
// 		return fmt.Errorf("wrong sell value")
// 	}
//
// 	c.ID = p.Id
// 	c.Buy = buy
// 	c.Sell = sell
//
// 	return nil
// }
