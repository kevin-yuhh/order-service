package errorm

import "errors"

var (
	AccountIllegal      = errors.New("illegal account")
	RowNotChanged       = errors.New("row has not changed")
	RequestParamEmpty   = errors.New("request param is empty")
	InsufficientBalance = errors.New("insufficient balance")
)
