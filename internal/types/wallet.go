package types

import (
	"fmt"
	"github.com/shopspring/decimal"
	proto "order_service/generated/wallet"
)

type Wallet struct {
	ID       string          `json:"id"`
	UserID   string          `json:"user_id"`
	Currency string          `json:"currency"`
	Value    decimal.Decimal `json:"value"`
	IsCrypto bool            `json:"is_crypto"`
}

func NewWallet(id, currency, value, userID string, isCrypto bool) *Wallet {
	val, _ := decimal.NewFromString(value)

	return &Wallet{
		ID:       id,
		Currency: currency,
		Value:    val,
		UserID:   userID,
		IsCrypto: isCrypto,
	}
}

func (w *Wallet) String() string {
	return fmt.Sprintf(
		"[Wallet ID: %s, Currency: %s, Value: %s, User ID: %s, Is Crypto: %t]",
		w.ID,
		w.Currency,
		w.Value,
		w.UserID,
		w.IsCrypto,
	)
}

func (w *Wallet) ToProto() *proto.Wallet {
	return &proto.Wallet{
		Id:       w.ID,
		Currency: w.Currency,
		Value:    w.Value.String(),
		UserId:   w.UserID,
		IsCrypto: w.IsCrypto,
	}
}

func (w *Wallet) FromProto(p *proto.Wallet) error {
	value, err := decimal.NewFromString(p.Value)
	if err != nil || value.LessThan(decimal.Zero) {
		return fmt.Errorf("wrong wallet value")
	}

	w.ID = p.Id
	w.Currency = p.Currency
	w.Value = value
	w.UserID = p.UserId
	w.IsCrypto = p.IsCrypto

	return nil
}
