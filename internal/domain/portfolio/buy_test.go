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

func TestCostBasis(t *testing.T) {
	tests := []struct {
		name         string
		quantity     decimal.Decimal
		pricePerUnit decimal.Decimal
		fee          decimal.Decimal
		expected     decimal.Decimal
		expectError  bool
	}{
		{
			name:         "standard case",
			quantity:     decimal.NewFromInt(2),
			pricePerUnit: decimal.NewFromFloat(testPrice),
			fee:          decimal.NewFromFloat(testFee),
			expected:     decimal.NewFromFloat(200002.50),
			expectError:  false,
		},
		{
			name:         "zero quantity returns error",
			quantity:     decimal.Zero,
			pricePerUnit: decimal.NewFromFloat(testPrice),
			fee:          decimal.NewFromFloat(testFee),
			expectError:  true,
		},
		{
			name:         "zero fee",
			quantity:     decimal.NewFromInt(1),
			pricePerUnit: decimal.NewFromFloat(testPrice),
			fee:          decimal.Zero,
			expected:     decimal.NewFromFloat(testPrice),
			expectError:  false,
		},
		{
			name:         "fractional quantities",
			quantity:     decimal.NewFromFloat(testOneSatoshi),
			pricePerUnit: decimal.NewFromFloat(testPrice),
			fee:          decimal.NewFromFloat(testFee),
			expected:     decimal.NewFromFloat(testFee + 0.001),
			expectError:  false,
		},
		{
			name:         "negative quantities",
			quantity:     decimal.NewFromInt(-2),
			pricePerUnit: decimal.NewFromFloat(testPrice),
			fee:          decimal.NewFromFloat(testFee),
			expectError:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buy := Buy{
				Quantity:     tt.quantity,
				PricePerUnit: money.NewUSD(tt.pricePerUnit),
				Fee:          money.NewUSD(tt.fee),
			}
			costBasis, err := buy.CostBasis()
			if tt.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.True(t, costBasis.Amount.Equal(tt.expected), "expected %s, got %s", tt.expected, costBasis.Amount)
		})
	}
}

func TestCostBasisPerUnit(t *testing.T) {
	tests := []struct {
		name         string
		quantity     decimal.Decimal
		pricePerUnit decimal.Decimal
		fee          decimal.Decimal
		expected     decimal.Decimal
		expectError  bool
	}{
		{
			name:         "standard case",
			quantity:     decimal.NewFromInt(2),
			pricePerUnit: decimal.NewFromFloat(testPrice),
			fee:          decimal.NewFromFloat(testFee),
			expected:     decimal.NewFromFloat(100001.25),
			expectError:  false,
		},
		{
			name:         "zero quantity returns error",
			quantity:     decimal.Zero,
			pricePerUnit: decimal.NewFromFloat(testPrice),
			fee:          decimal.NewFromFloat(testFee),
			expectError:  true,
		},
		{
			name:         "negative quantity returns error",
			quantity:     decimal.NewFromInt(-2),
			pricePerUnit: decimal.NewFromFloat(testPrice),
			fee:          decimal.NewFromFloat(testFee),
			expectError:  true,
		},
		{
			name:         "fractional quantities",
			quantity:     decimal.NewFromFloat(testOneSatoshi),
			pricePerUnit: decimal.NewFromFloat(testPrice),
			fee:          decimal.NewFromFloat(testOneSatoshi),
			expected:     decimal.NewFromFloat(100001.00),
			expectError:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buy := Buy{
				Quantity:     tt.quantity,
				PricePerUnit: money.NewUSD(tt.pricePerUnit),
				Fee:          money.NewUSD(tt.fee),
			}
			costBasisPerUnit, err := buy.CostBasisPerUnit()
			if tt.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.True(t, costBasisPerUnit.Amount.Equal(tt.expected), "expected %s, got %s", tt.expected, costBasisPerUnit.Amount)
		})
	}
}

func TestCreateTaxLot(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name        string
		quantity    decimal.Decimal
		expectError bool
	}{
		{
			name:        "creates tax lot from valid buy",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buyID := uuid.New()
			walletID := uuid.New()
			pricePerUnit := decimal.NewFromFloat(testPrice)
			fee := decimal.NewFromFloat(testFee)

			buy := Buy{
				ID:           buyID,
				WalletID:     walletID,
				Hash:         "tx_hash_123",
				Timestamp:    fixedTime,
				Asset:        "BTC",
				Quantity:     tt.quantity,
				PricePerUnit: money.NewUSD(pricePerUnit),
				Fee:          money.NewUSD(fee),
			}

			taxLot, err := buy.CreateTaxLot()

			if tt.expectError {
				require.Error(t, err)
				require.Equal(t, TaxLot{}, taxLot)
				return
			}

			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, taxLot.ID)
			assert.NotEqual(t, buyID, taxLot.ID)
			assert.Equal(t, walletID, taxLot.WalletID)
			assert.Equal(t, "BTC", taxLot.Asset)
			assert.Equal(t, fixedTime, taxLot.AcquisitionDate)
			assert.True(t, tt.quantity.Equal(taxLot.OriginalQuantity))
			assert.True(t, tt.quantity.Equal(taxLot.RemainingQuantity))
			expectedCostBasis, err := buy.CostBasis()
			require.NoError(t, err)
			assert.True(t, expectedCostBasis.Amount.Equal(taxLot.OriginalCostBasis.Amount))
			assert.Equal(t, buy, taxLot.SourceBuy)
		})
	}
}
