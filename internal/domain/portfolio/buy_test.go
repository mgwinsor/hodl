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
	}{
		{
			name:         "standard case",
			quantity:     decimal.NewFromInt(2),
			pricePerUnit: decimal.NewFromFloat(testPrice),
			fee:          decimal.NewFromFloat(testFee),
			expected:     decimal.NewFromFloat(200002.50),
		},
		{
			name:         "zero fee",
			quantity:     decimal.NewFromInt(1),
			pricePerUnit: decimal.NewFromFloat(testPrice),
			fee:          decimal.Zero,
			expected:     decimal.NewFromFloat(testPrice),
		},
		{
			name:         "fractional quantities",
			quantity:     decimal.NewFromFloat(testOneSatoshi),
			pricePerUnit: decimal.NewFromFloat(testPrice),
			fee:          decimal.NewFromFloat(testFee),
			expected:     decimal.NewFromFloat(testFee + 0.001),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qty, err := NewQuantity(tt.quantity)
			require.NoError(t, err)

			buy := Buy{
				Quantity:     qty,
				PricePerUnit: money.NewUSD(tt.pricePerUnit),
				Fee:          money.NewUSD(tt.fee),
			}
			costBasis := buy.CostBasis()
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
	}{
		{
			name:         "standard case",
			quantity:     decimal.NewFromInt(2),
			pricePerUnit: decimal.NewFromFloat(testPrice),
			fee:          decimal.NewFromFloat(testFee),
			expected:     decimal.NewFromFloat(100001.25),
		},
		{
			name:         "fractional quantities",
			quantity:     decimal.NewFromFloat(testOneSatoshi),
			pricePerUnit: decimal.NewFromFloat(testPrice),
			fee:          decimal.NewFromFloat(testOneSatoshi),
			expected:     decimal.NewFromFloat(100001.00),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qty, err := NewQuantity(tt.quantity)
			require.NoError(t, err)

			buy := Buy{
				Quantity:     qty,
				PricePerUnit: money.NewUSD(tt.pricePerUnit),
				Fee:          money.NewUSD(tt.fee),
			}
			costBasisPerUnit := buy.CostBasisPerUnit()
			require.True(t, costBasisPerUnit.Amount.Equal(tt.expected), "expected %s, got %s", tt.expected, costBasisPerUnit.Amount)
		})
	}
}

func TestCreateTaxLot(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name     string
		quantity decimal.Decimal
	}{
		{
			name:     "creates tax lot from valid buy",
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
			buyID := uuid.New()
			walletID := uuid.New()
			pricePerUnit := decimal.NewFromFloat(testPrice)
			fee := decimal.NewFromFloat(testFee)

			qty, err := NewQuantity(tt.quantity)
			require.NoError(t, err)

			buy := Buy{
				ID:           buyID,
				WalletID:     walletID,
				Hash:         "tx_hash_123",
				Timestamp:    fixedTime,
				Asset:        "BTC",
				Quantity:     qty,
				PricePerUnit: money.NewUSD(pricePerUnit),
				Fee:          money.NewUSD(fee),
			}

			taxLot := buy.CreateTaxLot()

			assert.NotEqual(t, uuid.Nil, taxLot.ID)
			assert.NotEqual(t, buyID, taxLot.ID)
			assert.Equal(t, walletID, taxLot.WalletID)
			assert.Equal(t, "BTC", taxLot.Asset)
			assert.Equal(t, fixedTime, taxLot.AcquisitionDate)
			assert.True(t, tt.quantity.Equal(taxLot.OriginalQuantity))
			assert.True(t, tt.quantity.Equal(taxLot.RemainingQuantity))
			expectedCostBasis := buy.CostBasis()
			assert.True(t, expectedCostBasis.Amount.Equal(taxLot.OriginalCostBasis.Amount))
			assert.Equal(t, buy, taxLot.Source)
		})
	}
}
