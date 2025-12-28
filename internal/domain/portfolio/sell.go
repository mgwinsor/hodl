package portfolio

import (
	"time"

	"github.com/google/uuid"
	"github.com/mgwinsor/hodl/internal/domain/money"
)

type Sell struct {
	ID           uuid.UUID
	WalletID     uuid.UUID
	Hash         string
	Timestamp    time.Time
	Asset        string
	Quantity     Quantity
	PricePerUnit money.USD
	Fee          money.USD
}

func (s Sell) GrossProceeds() money.USD {
	grossProceeds := s.Quantity.Value().Mul(s.PricePerUnit.Amount)
	return money.NewUSD(grossProceeds)
}

func (s Sell) NetProceeds() money.USD {
	gross := s.GrossProceeds()
	net := gross.Amount.Sub(s.Fee.Amount)
	return money.NewUSD(net)
}
