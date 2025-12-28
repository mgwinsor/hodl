package portfolio

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mgwinsor/hodl/internal/domain/money"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReward_TaxableIncome(t *testing.T) {
	tests := []struct {
		name            string
		quantity        decimal.Decimal
		fairMarketValue decimal.Decimal
		expected        decimal.Decimal
	}{
		{
			name:            "standard case",
			quantity:        decimal.NewFromInt(2),
			fairMarketValue: decimal.NewFromFloat(testPrice),
			expected:        decimal.NewFromFloat(testPrice),
		},
		{
			name:            "fractional quantity",
			quantity:        decimal.NewFromFloat(testOneSatoshi),
			fairMarketValue: decimal.NewFromFloat(0.001),
			expected:        decimal.NewFromFloat(0.001),
		},
		{
			name:            "large quantity",
			quantity:        decimal.NewFromInt(1000000),
			fairMarketValue: decimal.NewFromFloat(50000000.00),
			expected:        decimal.NewFromFloat(50000000.00),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qty, err := NewQuantity(tt.quantity)
			require.NoError(t, err)

			reward := Reward{
				Quantity:        qty,
				FairMarketValue: money.NewUSD(tt.fairMarketValue),
			}

			taxableIncome := reward.TaxableIncome()
			require.True(t, taxableIncome.Amount.Equal(tt.expected), "expected %s, got %s", tt.expected, taxableIncome.Amount)
		})
	}
}

func TestReward_CostBasisPerUnit(t *testing.T) {
	tests := []struct {
		name            string
		quantity        decimal.Decimal
		fairMarketValue decimal.Decimal
		expected        decimal.Decimal
	}{
		{
			name:            "standard case",
			quantity:        decimal.NewFromInt(2),
			fairMarketValue: decimal.NewFromFloat(200000.00),
			expected:        decimal.NewFromFloat(testPrice),
		},
		{
			name:            "fractional quantity",
			quantity:        decimal.NewFromFloat(testOneSatoshi),
			fairMarketValue: decimal.NewFromFloat(0.001),
			expected:        decimal.NewFromFloat(100000.00),
		},
		{
			name:            "single unit",
			quantity:        decimal.NewFromInt(1),
			fairMarketValue: decimal.NewFromFloat(testPrice),
			expected:        decimal.NewFromFloat(testPrice),
		},
		{
			name:            "non-even division",
			quantity:        decimal.NewFromInt(3),
			fairMarketValue: decimal.NewFromFloat(100.00),
			expected:        decimal.RequireFromString("33.3333333333333333"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qty, err := NewQuantity(tt.quantity)
			require.NoError(t, err)

			reward := Reward{
				Quantity:        qty,
				FairMarketValue: money.NewUSD(tt.fairMarketValue),
			}

			costBasisPerUnit := reward.CostBasisPerUnit()
			require.True(t, costBasisPerUnit.Amount.Equal(tt.expected), "expected %s, got %s", tt.expected, costBasisPerUnit.Amount)
		})
	}
}

func TestReward_CreateTaxLot(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name     string
		quantity decimal.Decimal
	}{
		{
			name:     "creates tax lot from valid reward",
			quantity: decimal.NewFromFloat(0.5),
		},
		{
			name:     "fractional quantity",
			quantity: decimal.NewFromFloat(testOneSatoshi),
		},
		{
			name:     "large quantity",
			quantity: decimal.NewFromInt(1000000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rewardID := uuid.New()
			walletID := uuid.New()
			fairMarketValue := decimal.NewFromFloat(testPrice)

			qty, err := NewQuantity(tt.quantity)
			require.NoError(t, err)

			reward := Reward{
				ID:              rewardID,
				WalletID:        walletID,
				Hash:            "reward_hash_456",
				Timestamp:       fixedTime,
				Asset:           "ETH",
				Quantity:        qty,
				FairMarketValue: money.NewUSD(fairMarketValue),
				Type:            RewardTypeStaking,
			}

			taxLot := reward.CreateTaxLot()

			assert.NotEqual(t, uuid.Nil, taxLot.ID)
			assert.NotEqual(t, rewardID, taxLot.ID)
			assert.Equal(t, walletID, taxLot.WalletID)
			assert.Equal(t, "ETH", taxLot.Asset)
			assert.Equal(t, fixedTime, taxLot.AcquisitionDate)
			assert.True(t, tt.quantity.Equal(taxLot.OriginalQuantity))
			assert.True(t, tt.quantity.Equal(taxLot.RemainingQuantity))
			assert.True(t, fairMarketValue.Equal(taxLot.OriginalCostBasis.Amount))
			assert.Equal(t, reward, taxLot.Source)
		})
	}
}
