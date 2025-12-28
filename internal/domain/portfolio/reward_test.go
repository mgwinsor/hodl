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
		expectError     bool
	}{
		{
			name:            "standard case",
			quantity:        decimal.NewFromInt(2),
			fairMarketValue: decimal.NewFromFloat(testPrice),
			expected:        decimal.NewFromFloat(testPrice),
			expectError:     false,
		},
		{
			name:            "zero quantity returns error",
			quantity:        decimal.Zero,
			fairMarketValue: decimal.NewFromFloat(testPrice),
			expectError:     true,
		},
		{
			name:            "negative quantity returns error",
			quantity:        decimal.NewFromInt(-1),
			fairMarketValue: decimal.NewFromFloat(testPrice),
			expectError:     true,
		},
		{
			name:            "fractional quantity",
			quantity:        decimal.NewFromFloat(testOneSatoshi),
			fairMarketValue: decimal.NewFromFloat(0.001),
			expected:        decimal.NewFromFloat(0.001),
			expectError:     false,
		},
		{
			name:            "large quantity",
			quantity:        decimal.NewFromInt(1000000),
			fairMarketValue: decimal.NewFromFloat(50000000.00),
			expected:        decimal.NewFromFloat(50000000.00),
			expectError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reward := Reward{
				Quantity:        tt.quantity,
				FairMarketValue: money.NewUSD(tt.fairMarketValue),
			}

			taxableIncome, err := reward.TaxableIncome()

			if tt.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
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
		expectError     bool
	}{
		{
			name:            "standard case",
			quantity:        decimal.NewFromInt(2),
			fairMarketValue: decimal.NewFromFloat(200000.00),
			expected:        decimal.NewFromFloat(testPrice),
			expectError:     false,
		},
		{
			name:            "zero quantity returns error",
			quantity:        decimal.Zero,
			fairMarketValue: decimal.NewFromFloat(testPrice),
			expectError:     true,
		},
		{
			name:            "negative quantity returns error",
			quantity:        decimal.NewFromInt(-1),
			fairMarketValue: decimal.NewFromFloat(testPrice),
			expectError:     true,
		},
		{
			name:            "fractional quantity",
			quantity:        decimal.NewFromFloat(testOneSatoshi),
			fairMarketValue: decimal.NewFromFloat(0.001),
			expected:        decimal.NewFromFloat(100000.00),
			expectError:     false,
		},
		{
			name:            "single unit",
			quantity:        decimal.NewFromInt(1),
			fairMarketValue: decimal.NewFromFloat(testPrice),
			expected:        decimal.NewFromFloat(testPrice),
			expectError:     false,
		},
		{
			name:            "non-even division",
			quantity:        decimal.NewFromInt(3),
			fairMarketValue: decimal.NewFromFloat(100.00),
			expected:        decimal.RequireFromString("33.3333333333333333"),
			expectError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reward := Reward{
				Quantity:        tt.quantity,
				FairMarketValue: money.NewUSD(tt.fairMarketValue),
			}

			costBasisPerUnit, err := reward.CostBasisPerUnit()

			if tt.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.True(t, costBasisPerUnit.Amount.Equal(tt.expected), "expected %s, got %s", tt.expected, costBasisPerUnit.Amount)
		})
	}
}

func TestReward_CreateTaxLot(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name        string
		quantity    decimal.Decimal
		expectError bool
	}{
		{
			name:        "creates tax lot from valid reward",
			quantity:    decimal.NewFromFloat(0.5),
			expectError: false,
		},
		{
			name:        "zero quantity returns error",
			quantity:    decimal.Zero,
			expectError: true,
		},
		{
			name:        "negative quantity returns error",
			quantity:    decimal.NewFromInt(-1),
			expectError: true,
		},
		{
			name:        "fractional quantity",
			quantity:    decimal.NewFromFloat(testOneSatoshi),
			expectError: false,
		},
		{
			name:        "large quantity",
			quantity:    decimal.NewFromInt(1000000),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rewardID := uuid.New()
			walletID := uuid.New()
			fairMarketValue := decimal.NewFromFloat(testPrice)

			reward := Reward{
				ID:              rewardID,
				WalletID:        walletID,
				Hash:            "reward_hash_456",
				Timestamp:       fixedTime,
				Asset:           "ETH",
				Quantity:        tt.quantity,
				FairMarketValue: money.NewUSD(fairMarketValue),
				Type:            RewardTypeStaking,
			}

			taxLot, err := reward.CreateTaxLot()

			if tt.expectError {
				require.Error(t, err)
				require.Equal(t, TaxLot{}, taxLot)
				return
			}

			require.NoError(t, err)
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
