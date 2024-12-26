package bussiness_errors

import "fmt"

type ErrorCode int

const (
	ErrInsufficientFiatCurrency ErrorCode = iota + 1
	ErrInsufficientCryptoCurrency
	ErrOrderAlreadyExists
	ErrOrderNotFound
	ErrPriceNotFound
)

const (
	MsgInsufficientFiatCurrency   = "Insufficient fiat currency"
	MsgInsufficientCryptoCurrency = "Insufficient crypto currency"
	MsgOrderAlreadyExists         = "Order already exists"
	MsgOrderNotFound              = "Order not found"
	MsgPriceNotFound              = "Price not found"
)

type CustomError struct {
	Code    ErrorCode
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

func NewCustomError(code ErrorCode, message string) error {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}
