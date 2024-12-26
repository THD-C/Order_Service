package types

type PendingOrder struct {
	Order          *Order
	FiatCurrency   string
	CryptoCurrency string
}
