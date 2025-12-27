package portfolio

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mgwinsor/hodl/internal/domain/money"
	"github.com/shopspring/decimal"
)

type Sell struct {
	ID           uuid.UUID
	WalletID     uuid.UUID
	Hash         string
	Timestamp    time.Time
	Asset        string
	Quantity     decimal.Decimal
	PricePerUnit money.USD
	Fee          money.USD
}

func (s Sell) GrossProceeds() (money.USD, error) {
	if s.Quantity.IsZero() {
		return money.USD{}, errors.New(
			"cannot calculate proceeds for zero quantity",
		)
	}
	grossProceeds := s.Quantity.Mul(s.PricePerUnit.Amount)
	return money.NewUSD(grossProceeds), nil
}

func (s Sell) NetProceeds() (money.USD, error) {
	gross, err := s.GrossProceeds()
	if err != nil {
		return money.USD{}, err
	}
	net := gross.Amount.Sub(s.Fee.Amount)
	return money.NewUSD(net), nil
}
