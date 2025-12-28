package portfolio

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mgwinsor/hodl/internal/domain/money"
	"github.com/shopspring/decimal"
)

type Buy struct {
	ID           uuid.UUID
	WalletID     uuid.UUID
	Hash         string
	Timestamp    time.Time
	Asset        string
	Quantity     decimal.Decimal
	PricePerUnit money.USD
	Fee          money.USD
}

func (b Buy) CostBasis() (money.USD, error) {
	if b.Quantity.IsZero() {
		return money.USD{}, errors.New(
			"cannot calculate cost basis for zero quantity",
		)
	}
	if b.Quantity.IsNegative() {
		return money.USD{}, errors.New(
			"cannot calculate cost basis for negative quantity",
		)
	}
	costBasis := b.Quantity.Mul(b.PricePerUnit.Amount).Add(b.Fee.Amount)
	return money.NewUSD(costBasis), nil
}

func (b Buy) CostBasisPerUnit() (money.USD, error) {
	totalCost, err := b.CostBasis()
	if err != nil {
		return money.USD{}, err
	}
	perUnit := totalCost.Amount.Div(b.Quantity)
	return money.NewUSD(perUnit), nil
}

func (b Buy) CreateTaxLot() (TaxLot, error) {
	costBasis, err := b.CostBasis()
	if err != nil {
		return TaxLot{}, err
	}

	return TaxLot{
		ID:                uuid.New(),
		WalletID:          b.WalletID,
		Asset:             b.Asset,
		AcquisitionDate:   b.Timestamp,
		OriginalQuantity:  b.Quantity,
		RemainingQuantity: b.Quantity,
		OriginalCostBasis: costBasis,
		Source:            b,
	}, nil
}
