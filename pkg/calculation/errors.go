package calculation

import "errors"

var (
	ErrInvalidExpression = errors.New("invalid expression")
	ErrDivisionByZero    = errors.New("division by zero")
	ErrUnknown           = errors.New("internal server error")
	ErrStrangeSymbols    = errors.New("symbols not like given")
)
