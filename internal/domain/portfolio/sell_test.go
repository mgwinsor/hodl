package portfolio

import (
	"testing"

	"github.com/mgwinsor/hodl/internal/domain/money"
	"github.com/shopspring/decimal"
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
			pricePerUnit: decimal.NewFromFloat(testPrice),
			expected:     decimal.NewFromFloat(200000.00),
			expectError:  false,
		},
		{
			name:         "zero quantity returns error",
			quantity:     decimal.Zero,
			pricePerUnit: decimal.NewFromFloat(testPrice),
			expectError:  true,
		},
		{
			name:         "fractional quantities",
			quantity:     decimal.NewFromFloat(testOneSatoshi),
			pricePerUnit: decimal.NewFromFloat(testPrice),
			expected:     decimal.NewFromFloat(0.001),
			expectError:  false,
		},
		{
			name:         "negative quantity",
			quantity:     decimal.NewFromInt(-2),
			pricePerUnit: decimal.NewFromFloat(testPrice),
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
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
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
			pricePerUnit: decimal.NewFromFloat(testPrice),
			fee:          decimal.NewFromFloat(testFee),
			expected:     decimal.NewFromFloat(199997.5),
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
			fee:          decimal.NewFromFloat(0.00001),
			expected:     decimal.NewFromFloat(0.00099),
			expectError:  false,
		},
		{
			name:         "fee exceeds gross proceeds results in negative net",
			quantity:     decimal.NewFromInt(1),
			pricePerUnit: decimal.NewFromFloat(100.00),
			fee:          decimal.NewFromFloat(150.00),
			expected:     decimal.NewFromFloat(-50.00),
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
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.True(t, netProceeds.Amount.Equal(tt.expected), "expected %s, got %s", tt.expected, netProceeds.Amount)
		})
	}
}
