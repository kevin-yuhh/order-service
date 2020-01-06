package errorm

import "errors"

var (
	AccountIllegal      = errors.New("illegal account")
	RowNotChanged       = errors.New("row has not changed")
	RequestParamEmpty   = errors.New("request param is empty")
	InsufficientBalance = errors.New("insufficient balance")
	OrderStatusIllegal  = errors.New("order status illegal")
	OrderTypeIllegal    = errors.New("order type illegal")
	OrderNotExists      = errors.New("order not exists")
	BtfsStatusNotExists = errors.New("btfs status not exists")
)
