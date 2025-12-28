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
	}{
		{
			name:         "standard case",
			quantity:     decimal.NewFromInt(2),
			pricePerUnit: decimal.NewFromFloat(testPrice),
			expected:     decimal.NewFromFloat(200000.00),
		},
		{
			name:         "fractional quantities",
			quantity:     decimal.NewFromFloat(testOneSatoshi),
			pricePerUnit: decimal.NewFromFloat(testPrice),
			expected:     decimal.NewFromFloat(0.001),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qty, err := NewQuantity(tt.quantity)
			require.NoError(t, err)

			sell := Sell{
				Quantity:     qty,
				PricePerUnit: money.NewUSD(tt.pricePerUnit),
			}
			grossProceeds := sell.GrossProceeds()
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
	}{
		{
			name:         "standard case",
			quantity:     decimal.NewFromInt(2),
			pricePerUnit: decimal.NewFromFloat(testPrice),
			fee:          decimal.NewFromFloat(testFee),
			expected:     decimal.NewFromFloat(199997.5),
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
			fee:          decimal.NewFromFloat(0.00001),
			expected:     decimal.NewFromFloat(0.00099),
		},
		{
			name:         "fee exceeds gross proceeds results in negative net",
			quantity:     decimal.NewFromInt(1),
			pricePerUnit: decimal.NewFromFloat(100.00),
			fee:          decimal.NewFromFloat(150.00),
			expected:     decimal.NewFromFloat(-50.00),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qty, err := NewQuantity(tt.quantity)
			require.NoError(t, err)

			sell := Sell{
				Quantity:     qty,
				PricePerUnit: money.NewUSD(tt.pricePerUnit),
				Fee:          money.NewUSD(tt.fee),
			}
			netProceeds := sell.NetProceeds()
			require.True(t, netProceeds.Amount.Equal(tt.expected), "expected %s, got %s", tt.expected, netProceeds.Amount)
		})
	}
}
