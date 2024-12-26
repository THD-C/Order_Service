package types

import (
	"fmt"
	"github.com/shopspring/decimal"
)

type Price struct {
	Bid decimal.Decimal
	Ask decimal.Decimal
}

func (p *Price) String() string {
	return fmt.Sprintf("Bid: %s, Ask: %s", p.Bid.String(), p.Ask.String())
}
