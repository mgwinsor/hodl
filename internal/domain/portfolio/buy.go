package portfolio

import (
	"time"

	"github.com/google/uuid"
	"github.com/mgwinsor/hodl/internal/domain/money"
)

type Buy struct {
	ID           uuid.UUID
	WalletID     uuid.UUID
	Hash         string
	Timestamp    time.Time
	Asset        string
	Quantity     Quantity
	PricePerUnit money.USD
	Fee          money.USD
}

func (b Buy) CostBasis() money.USD {
	costBasis := b.Quantity.Value().Mul(b.PricePerUnit.Amount).Add(b.Fee.Amount)
	return money.NewUSD(costBasis)
}

func (b Buy) CostBasisPerUnit() money.USD {
	totalCost := b.CostBasis()
	perUnit := totalCost.Amount.Div(b.Quantity.Value())
	return money.NewUSD(perUnit)
}

func (b Buy) CreateTaxLot() TaxLot {
	return TaxLot{
		ID:                uuid.New(),
		WalletID:          b.WalletID,
		Asset:             b.Asset,
		AcquisitionDate:   b.Timestamp,
		OriginalQuantity:  b.Quantity.Value(),
		RemainingQuantity: b.Quantity.Value(),
		OriginalCostBasis: b.CostBasis(),
		Source:            b,
	}
}
