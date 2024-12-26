package types

import (
	"fmt"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"
	proto "order_service/generated/order"
)

type Order struct {
	ID             string                 `json:"id"`
	UserID         string                 `json:"user_id"`
	DateCreated    *timestamppb.Timestamp `json:"date_created"`
	DateExecuted   *timestamppb.Timestamp `json:"date_executed"`
	Status         proto.OrderStatus      `json:"status"`
	Nominal        decimal.Decimal        `json:"nominal"`
	CashQuantity   decimal.Decimal        `json:"cash_quantity"`
	Price          decimal.Decimal        `json:"price"`
	Type           proto.OrderType        `json:"type"`
	Side           proto.OrderSide        `json:"side"`
	CryptoWalletID string                 `json:"crypto_wallet_id"`
	FiatWalletID   string                 `json:"fiat_wallet_id"`
}

func NewOrder(
	id, userID, nominal, cashQuantity, price, cryptoWalletID, fiatWalletID string,
	dateCreated, dateExecuted *timestamppb.Timestamp,
	status proto.OrderStatus,
	orderType proto.OrderType,
	side proto.OrderSide,
) *Order {
	nom, _ := decimal.NewFromString(nominal)
	cashQty, _ := decimal.NewFromString(cashQuantity)
	prc, _ := decimal.NewFromString(price)

	return &Order{
		ID:             id,
		UserID:         userID,
		DateCreated:    dateCreated,
		DateExecuted:   dateExecuted,
		Status:         status,
		Nominal:        nom,
		CashQuantity:   cashQty,
		Price:          prc,
		Type:           orderType,
		Side:           side,
		CryptoWalletID: cryptoWalletID,
		FiatWalletID:   fiatWalletID,
	}
}

func (o *Order) String() string {
	return fmt.Sprintf(
		"[Order ID: %s User ID: %s, Date Created: %s, Date Executed: %s, Status: %s, Nominal: %s"+
			", Cash Quantity: %s, Price: %s, Type: %s, Side: %s, "+
			"Crypto Wallet ID: %s, Fiat Wallet ID: %s]",
		o.ID,
		o.UserID,
		o.DateCreated,
		o.DateExecuted,
		o.Status,
		o.Nominal,
		o.CashQuantity,
		o.Price,
		o.Type,
		o.Side,
		o.CryptoWalletID,
		o.FiatWalletID,
	)
}

func (o *Order) ToProto() *proto.OrderDetails {
	return &proto.OrderDetails{
		Id:             o.ID,
		UserId:         o.UserID,
		DateCreated:    o.DateCreated,
		DateExecuted:   o.DateExecuted,
		Status:         o.Status,
		Nominal:        o.Nominal.String(),
		CashQuantity:   o.CashQuantity.String(),
		Price:          o.Price.String(),
		Type:           o.Type,
		Side:           o.Side,
		CryptoWalletId: o.CryptoWalletID,
		FiatWalletId:   o.FiatWalletID,
	}
}

func (o *Order) FromProto(p *proto.OrderDetails) error {
	nominal, err := decimal.NewFromString(p.Nominal)
	if err != nil {
		return err
	}

	cashQuantity, err := decimal.NewFromString(p.CashQuantity)
	if err != nil {
		return err
	}

	price, err := decimal.NewFromString(p.Price)
	if err != nil {
		return err
	}

	o.ID = p.Id
	o.UserID = p.UserId
	o.DateCreated = p.DateCreated
	o.DateExecuted = p.DateExecuted
	o.Status = p.Status
	o.Nominal = nominal
	o.CashQuantity = cashQuantity
	o.Price = price
	o.Type = p.Type
	o.Side = p.Side
	o.CryptoWalletID = p.CryptoWalletId
	o.FiatWalletID = p.FiatWalletId

	return nil
}
