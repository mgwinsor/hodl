package portfolio

import (
	"testing"

	"github.com/mgwinsor/hodl/internal/domain/money"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrossProceeds(t *testing.T) {
	tests := []struct {
		name         string
		quantity     decimal.Decimal
		pricePerUnit decimal.Decimal
		expected     decimal.Decimal
		expectError  bool
	}{
		{
			name:         "standard case",
			quantity:     decimal.NewFromInt(2),
			pricePerUnit: decimal.NewFromFloat(80000.00),
			expected:     decimal.NewFromFloat(160000.00),
			expectError:  false,
		},
		{
			name:         "zero quantity returns error",
			quantity:     decimal.Zero,
			pricePerUnit: decimal.NewFromFloat(80000.00),
			expectError:  true,
		},
		{
			name:         "fractional quantities",
			quantity:     decimal.NewFromFloat(0.001),
			pricePerUnit: decimal.NewFromFloat(80000.00),
			expected:     decimal.NewFromFloat(80.00),
			expectError:  false,
		},
		{
			name:         "negative quantity",
			quantity:     decimal.NewFromInt(-2),
			pricePerUnit: decimal.NewFromFloat(80000.00),
			expectError:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sell := Sell{
				Quantity:     tt.quantity,
				PricePerUnit: money.NewUSD(tt.pricePerUnit),
			}
			grossProceeds, err := sell.GrossProceeds()
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			require.True(t, grossProceeds.Amount.Equal(tt.expected), "expected %s, got %s", tt.expected, grossProceeds.Amount)
		})
	}
}

func TestNetProceeds(t *testing.T) {
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
			pricePerUnit: decimal.NewFromFloat(80000.00),
			fee:          decimal.NewFromFloat(2.50),
			expected:     decimal.NewFromFloat(159997.5),
			expectError:  false,
		},
		{
			name:         "zero quantity returns error",
			quantity:     decimal.Zero,
			pricePerUnit: decimal.NewFromFloat(80000.00),
			fee:          decimal.NewFromFloat(2.50),
			expectError:  true,
		},
		{
			name:         "zero fee",
			quantity:     decimal.NewFromInt(1),
			pricePerUnit: decimal.NewFromFloat(80000.00),
			fee:          decimal.Zero,
			expected:     decimal.NewFromFloat(80000.00),
			expectError:  false,
		},
		{
			name:         "fractional quantities",
			quantity:     decimal.NewFromFloat(0.001),
			pricePerUnit: decimal.NewFromFloat(80000.00),
			fee:          decimal.NewFromFloat(0.99),
			expected:     decimal.NewFromFloat(79.01),
			expectError:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sell := Sell{
				Quantity:     tt.quantity,
				PricePerUnit: money.NewUSD(tt.pricePerUnit),
				Fee:          money.NewUSD(tt.fee),
			}
			netProceeds, err := sell.NetProceeds()
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			require.True(t, netProceeds.Amount.Equal(tt.expected), "expected %s, got %s", tt.expected, netProceeds.Amount)
		})
	}
}
