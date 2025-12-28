package portfolio

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mgwinsor/hodl/internal/domain/money"
	"github.com/shopspring/decimal"
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
	Quantity        decimal.Decimal
	FairMarketValue money.USD
	Type            RewardType
}

func (r Reward) TaxableIncome() (money.USD, error) {
	if r.Quantity.IsZero() {
		return money.USD{}, errors.New(
			"cannot calculate taxable income for zero quantity",
		)
	}
	if r.Quantity.IsNegative() {
		return money.USD{}, errors.New(
			"cannot calculate taxable income for negative quantity",
		)
	}
	return r.FairMarketValue, nil
}

func (r Reward) CostBasisPerUnit() (money.USD, error) {
	if r.Quantity.IsZero() {
		return money.USD{}, errors.New(
			"cannot calculate cost basis per unit for zero quantity",
		)
	}
	if r.Quantity.IsNegative() {
		return money.USD{}, errors.New(
			"cannot calculate cost basis per unit for negative quantity",
		)
	}
	perUnit := r.FairMarketValue.Amount.Div(r.Quantity)
	return money.NewUSD(perUnit), nil
}

func (r Reward) CreateTaxLot() (TaxLot, error) {
	if r.Quantity.IsZero() {
		return TaxLot{}, errors.New(
			"cannot create tax lot for zero quantity",
		)
	}
	if r.Quantity.IsNegative() {
		return TaxLot{}, errors.New(
			"cannot create tax lot for negative quantity",
		)
	}

	return TaxLot{
		ID:                uuid.New(),
		WalletID:          r.WalletID,
		Asset:             r.Asset,
		AcquisitionDate:   r.Timestamp,
		OriginalQuantity:  r.Quantity,
		RemainingQuantity: r.Quantity,
		OriginalCostBasis: r.FairMarketValue,
		Source:            r,
	}, nil
}
