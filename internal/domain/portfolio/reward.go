package portfolio

import (
	"time"

	"github.com/google/uuid"
	"github.com/mgwinsor/hodl/internal/domain/money"
)

type RewardType string

const (
	RewardTypeStaking  RewardType = "staking"
	RewardTypeAirdrop  RewardType = "airdrop"
	RewardTypeMining   RewardType = "mining"
	RewardTypeInterest RewardType = "interest"
)

type Reward struct {
	ID              uuid.UUID
	WalletID        uuid.UUID
	Hash            string
	Timestamp       time.Time
	Asset           string
	Quantity        Quantity
	FairMarketValue money.USD
	Type            RewardType
}

func (r Reward) TaxableIncome() money.USD {
	return r.FairMarketValue
}

func (r Reward) CostBasisPerUnit() money.USD {
	perUnit := r.FairMarketValue.Amount.Div(r.Quantity.Value())
	return money.NewUSD(perUnit)
}

func (r Reward) CreateTaxLot() TaxLot {
	return TaxLot{
		ID:                uuid.New(),
		WalletID:          r.WalletID,
		Asset:             r.Asset,
		AcquisitionDate:   r.Timestamp,
		OriginalQuantity:  r.Quantity.Value(),
		RemainingQuantity: r.Quantity.Value(),
		OriginalCostBasis: r.FairMarketValue,
		Source:            r,
	}
}
