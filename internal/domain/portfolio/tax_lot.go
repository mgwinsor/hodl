package portfolio

import (
	"time"

	"github.com/google/uuid"
	"github.com/mgwinsor/hodl/internal/domain/money"
	"github.com/shopspring/decimal"
)

type TaxLotSource interface {
	CostBasisPerUnit() money.USD
}

type TaxLot struct {
	ID                uuid.UUID
	WalletID          uuid.UUID
	Asset             string
	AcquisitionDate   time.Time
	OriginalQuantity  decimal.Decimal
	RemainingQuantity decimal.Decimal
	OriginalCostBasis money.USD
	Source            TaxLotSource
}

func (lot TaxLot) IsFullyConsumed() bool {
	return lot.RemainingQuantity.IsZero()
}

func (lot TaxLot) HoldingPeriod(asOf time.Time) time.Duration {
	return asOf.Sub(lot.AcquisitionDate)
}

func (lot TaxLot) IsLongTerm(asOf time.Time) bool {
	oneYearLater := lot.AcquisitionDate.AddDate(1, 0, 0)
	return asOf.After(oneYearLater)
}
