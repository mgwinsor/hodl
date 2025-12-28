package portfolio

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewQuantity(t *testing.T) {
	tests := []struct {
		name        string
		value       decimal.Decimal
		expectError bool
	}{
		{
			name:        "positive integer",
			value:       decimal.NewFromInt(1),
			expectError: false,
		},
		{
			name:        "positive decimal",
			value:       decimal.NewFromFloat(0.5),
			expectError: false,
		},
		{
			name:        "very small positive value",
			value:       decimal.NewFromFloat(testOneSatoshi),
			expectError: false,
		},
		{
			name:        "large positive value",
			value:       decimal.NewFromInt(1000000),
			expectError: false,
		},
		{
			name:        "zero returns error",
			value:       decimal.Zero,
			expectError: true,
		},
		{
			name:        "negative integer returns error",
			value:       decimal.NewFromInt(-1),
			expectError: true,
		},
		{
			name:        "negative decimal returns error",
			value:       decimal.NewFromFloat(-0.5),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qty, err := NewQuantity(tt.value)

			if tt.expectError {
				require.Error(t, err)
				assert.Equal(t, Quantity{}, qty)
				return
			}

			require.NoError(t, err)
			assert.True(t, tt.value.Equal(qty.Value()))
		})
	}
}

func TestQuantity_Value(t *testing.T) {
	expected := decimal.NewFromFloat(1.5)
	qty, err := NewQuantity(expected)
	require.NoError(t, err)

	assert.True(t, expected.Equal(qty.Value()))
}
