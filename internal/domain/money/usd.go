package money

import "github.com/shopspring/decimal"

type USD struct {
	Amount decimal.Decimal
}

func NewUSD(amount decimal.Decimal) USD {
	return USD{Amount: amount}
}
