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

func TestFeeAsSell(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	transferID := uuid.New()
	fromWalletID := uuid.New()
	toWalletID := uuid.New()

	transferQty, err := NewQuantity(decimal.NewFromFloat(1.5))
	require.NoError(t, err)

	tests := []struct {
		name                 string
		feeQuantity          decimal.Decimal
		feeUSD               decimal.Decimal
		expectError          bool
		expectedPricePerUnit decimal.Decimal
	}{
		{
			name:                 "standard case with valid fee",
			feeQuantity:          decimal.NewFromFloat(0.001),
			feeUSD:               decimal.NewFromFloat(testPrice),
			expectError:          false,
			expectedPricePerUnit: decimal.NewFromFloat(100000000.00),
		},
		{
			name:        "zero fee quantity returns error",
			feeQuantity: decimal.Zero,
			feeUSD:      decimal.NewFromFloat(testFee),
			expectError: true,
		},
		{
			name:        "negative fee quantity returns error",
			feeQuantity: decimal.NewFromInt(-1),
			feeUSD:      decimal.NewFromFloat(testFee),
			expectError: true,
		},
		{
			name:                 "fractional fee quantity like satoshis",
			feeQuantity:          decimal.NewFromFloat(testOneSatoshi),
			feeUSD:               decimal.NewFromFloat(0.001),
			expectError:          false,
			expectedPricePerUnit: decimal.NewFromFloat(testPrice),
		},
		{
			name:                 "whole number fee quantity",
			feeQuantity:          decimal.NewFromInt(2),
			feeUSD:               decimal.NewFromFloat(200000.00),
			expectError:          false,
			expectedPricePerUnit: decimal.NewFromFloat(testPrice),
		},
		{
			name:                 "very small fee USD amount",
			feeQuantity:          decimal.NewFromFloat(testOneSatoshi),
			feeUSD:               decimal.NewFromFloat(testOneSatoshi),
			expectError:          false,
			expectedPricePerUnit: decimal.NewFromFloat(1.0),
		},
		{
			name:                 "non-integer price per unit result",
			feeQuantity:          decimal.NewFromInt(3),
			feeUSD:               decimal.NewFromFloat(100.00),
			expectError:          false,
			expectedPricePerUnit: decimal.NewFromFloat(100.00).Div(decimal.NewFromInt(3)),
		},
		{
			name:                 "very small quantity with large USD",
			feeQuantity:          decimal.NewFromFloat(testOneSatoshi),
			feeUSD:               decimal.NewFromFloat(testPrice),
			expectError:          false,
			expectedPricePerUnit: decimal.NewFromFloat(testPrice).Div(decimal.NewFromFloat(testOneSatoshi)),
		},
		{
			name:                 "quantity of one returns exact USD as price",
			feeQuantity:          decimal.NewFromInt(1),
			feeUSD:               decimal.NewFromFloat(testPrice),
			expectError:          false,
			expectedPricePerUnit: decimal.NewFromFloat(testPrice),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transfer := Transfer{
				ID:           transferID,
				Hash:         "tx_hash_transfer_123",
				Timestamp:    fixedTime,
				Asset:        "BTC",
				Quantity:     transferQty,
				FromWalletID: fromWalletID,
				ToWalletID:   toWalletID,
				FeeQuantity:  tt.feeQuantity,
				FeeAsset:     "ETH",
				FeeUSD:       money.NewUSD(tt.feeUSD),
			}

			sell, err := transfer.FeeAsSell()

			if tt.expectError {
				require.Error(t, err)
				require.Equal(t, Sell{}, sell)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, transferID, sell.ID)
			assert.Equal(t, fromWalletID, sell.WalletID)
			assert.Equal(t, "tx_hash_transfer_123", sell.Hash)
			assert.Equal(t, fixedTime, sell.Timestamp)
			assert.Equal(t, "ETH", sell.Asset)
			assert.True(t, tt.feeQuantity.Equal(sell.Quantity.Value()), "expected quantity %s, got %s", tt.feeQuantity, sell.Quantity.Value())
			assert.True(t, tt.expectedPricePerUnit.Equal(sell.PricePerUnit.Amount), "expected price per unit %s, got %s", tt.expectedPricePerUnit, sell.PricePerUnit.Amount)
			assert.True(t, decimal.Zero.Equal(sell.Fee.Amount), "expected zero fee, got %s", sell.Fee.Amount)
		})
	}
}
