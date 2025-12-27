package portfolio

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestIsFullyConsumed(t *testing.T) {
	tests := []struct {
		name              string
		remainingQuantity decimal.Decimal
		expected          bool
	}{
		{
			name:              "zero remaining quantity is fully consumed",
			remainingQuantity: decimal.Zero,
			expected:          true,
		},
		{
			name:              "positive remaining quantity is not fully consumed",
			remainingQuantity: decimal.NewFromInt(1),
			expected:          false,
		},
		{
			name:              "fractional remaining quantity is not fully consumed",
			remainingQuantity: decimal.NewFromFloat(testOneSatoshi),
			expected:          false,
		},
		{
			name:              "large remaining quantity is not fully consumed",
			remainingQuantity: decimal.NewFromInt(1000000),
			expected:          false,
		},
		{
			name:              "very small remaining quantity is not fully consumed",
			remainingQuantity: decimal.NewFromFloat(0.000000001),
			expected:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lot := TaxLot{
				RemainingQuantity: tt.remainingQuantity,
			}

			result := lot.IsFullyConsumed()

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHoldingPeriod(t *testing.T) {
	tests := []struct {
		name            string
		acquisitionDate time.Time
		asOf            time.Time
		expected        time.Duration
	}{
		{
			name:            "one day holding period",
			acquisitionDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			asOf:            time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			expected:        24 * time.Hour,
		},
		{
			name:            "exactly one year holding period non-leap year",
			acquisitionDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			asOf:            time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected:        365 * 24 * time.Hour,
		},
		{
			name:            "exactly one year holding period spanning leap year",
			acquisitionDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			asOf:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected:        366 * 24 * time.Hour,
		},
		{
			name:            "zero holding period same instant",
			acquisitionDate: time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC),
			asOf:            time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC),
			expected:        0,
		},
		{
			name:            "holding period with hours and minutes",
			acquisitionDate: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			asOf:            time.Date(2024, 1, 1, 14, 30, 0, 0, time.UTC),
			expected:        4*time.Hour + 30*time.Minute,
		},
		{
			name:            "leap year acquisition on Feb 29",
			acquisitionDate: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			asOf:            time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			expected:        24 * time.Hour,
		},
		{
			name:            "leap year Feb 29 to next year Feb 28",
			acquisitionDate: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			asOf:            time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
			expected:        365 * 24 * time.Hour,
		},
		{
			name:            "leap year Feb 29 to next year Mar 1",
			acquisitionDate: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			asOf:            time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			expected:        366 * 24 * time.Hour,
		},
		{
			name:            "multi-year holding period",
			acquisitionDate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			asOf:            time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected:        (365*4 + 1) * 24 * time.Hour, // 4 years including one leap year (2020)
		},
		{
			name:            "negative duration when asOf is before acquisition",
			acquisitionDate: time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
			asOf:            time.Date(2024, 6, 14, 0, 0, 0, 0, time.UTC),
			expected:        -24 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lot := TaxLot{
				AcquisitionDate: tt.acquisitionDate,
			}

			result := lot.HoldingPeriod(tt.asOf)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsLongTerm(t *testing.T) {
	tests := []struct {
		name            string
		acquisitionDate time.Time
		asOf            time.Time
		expected        bool
	}{
		{
			name:            "exactly one year is not long term",
			acquisitionDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
			asOf:            time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
			expected:        false,
		},
		{
			name:            "one year plus one day is long term",
			acquisitionDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
			asOf:            time.Date(2024, 6, 16, 0, 0, 0, 0, time.UTC),
			expected:        true,
		},
		{
			name:            "one year plus one second is long term",
			acquisitionDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
			asOf:            time.Date(2024, 6, 15, 0, 0, 1, 0, time.UTC),
			expected:        true,
		},
		{
			name:            "one year minus one second is not long term",
			acquisitionDate: time.Date(2023, 6, 15, 0, 0, 1, 0, time.UTC),
			asOf:            time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
			expected:        false,
		},
		{
			name:            "one day is not long term",
			acquisitionDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			asOf:            time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			expected:        false,
		},
		{
			name:            "same day is not long term",
			acquisitionDate: time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC),
			asOf:            time.Date(2024, 6, 15, 14, 0, 0, 0, time.UTC),
			expected:        false,
		},
		{
			name:            "multi-year holding is long term",
			acquisitionDate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			asOf:            time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected:        true,
		},
		{
			name:            "leap year Feb 29 exactly one year later to Feb 28 is not long term",
			acquisitionDate: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			asOf:            time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
			expected:        false, // AddDate(1,0,0) on Feb 29 2024 = Mar 1 2025
		},
		{
			name:            "leap year Feb 29 to next year Mar 1 is not long term",
			acquisitionDate: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			asOf:            time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			expected:        false, // exactly one year per AddDate
		},
		{
			name:            "leap year Feb 29 to next year Mar 2 is long term",
			acquisitionDate: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			asOf:            time.Date(2025, 3, 2, 0, 0, 0, 0, time.UTC),
			expected:        true,
		},
		{
			name:            "Feb 28 non-leap to Feb 28 next year is not long term",
			acquisitionDate: time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC),
			asOf:            time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC),
			expected:        false,
		},
		{
			name:            "Feb 28 non-leap to Feb 29 next leap year is long term",
			acquisitionDate: time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC),
			asOf:            time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			expected:        true,
		},
		{
			name:            "Dec 31 to Jan 1 next year plus one is long term",
			acquisitionDate: time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
			asOf:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected:        true,
		},
		{
			name:            "year boundary exactly one year",
			acquisitionDate: time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
			asOf:            time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			expected:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lot := TaxLot{
				AcquisitionDate: tt.acquisitionDate,
			}

			result := lot.IsLongTerm(tt.asOf)

			assert.Equal(t, tt.expected, result)
		})
	}
}
