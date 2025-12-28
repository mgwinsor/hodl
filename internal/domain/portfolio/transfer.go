package portfolio

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mgwinsor/hodl/internal/domain/money"
	"github.com/shopspring/decimal"
)

type Transfer struct {
	ID           uuid.UUID
	Hash         string
	Timestamp    time.Time
	Asset        string
	Quantity     Quantity
	FromWalletID uuid.UUID
	ToWalletID   uuid.UUID
	FeeQuantity  decimal.Decimal
	FeeAsset     string
	FeeUSD       money.USD
}

func (t Transfer) FeeAsSell() (Sell, error) {
	feeQty, err := NewQuantity(t.FeeQuantity)
	if err != nil {
		return Sell{}, errors.New(
			"cannot create sell for invalid fee quantity: " + err.Error(),
		)
	}
	pricePerUnit := t.FeeUSD.Amount.Div(t.FeeQuantity)
	return Sell{
		ID:           t.ID,
		WalletID:     t.FromWalletID,
		Hash:         t.Hash,
		Timestamp:    t.Timestamp,
		Asset:        t.FeeAsset,
		Quantity:     feeQty,
		PricePerUnit: money.NewUSD(pricePerUnit),
		Fee:          money.NewUSD(decimal.Zero),
	}, nil
}
