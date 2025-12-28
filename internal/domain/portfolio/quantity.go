package portfolio

import (
	"errors"

	"github.com/shopspring/decimal"
)

type Quantity struct {
	value decimal.Decimal
}

func NewQuantity(d decimal.Decimal) (Quantity, error) {
	if d.IsZero() {
		return Quantity{}, errors.New("quantity must be positive, got zero")
	}
	if d.IsNegative() {
		return Quantity{}, errors.New("quantity must be positive, got negative value")
	}
	return Quantity{value: d}, nil
}

func (q Quantity) Value() decimal.Decimal {
	return q.value
}
