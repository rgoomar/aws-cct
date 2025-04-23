package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceCosts_Initialization(t *testing.T) {
	// Test initialization of ServiceCosts
	costs := ServiceCosts{
		ServiceName:  "Amazon EC2",
		Amount:       100.00,
		SecondAmount: 120.00,
		Delta:        20.00,
		DeltaPercent: 20.00,
	}

	// Assertions
	assert.Equal(t, "Amazon EC2", costs.ServiceName)
	assert.Equal(t, 100.00, costs.Amount)
	assert.Equal(t, 120.00, costs.SecondAmount)
	assert.Equal(t, 20.00, costs.Delta)
	assert.Equal(t, 20.00, costs.DeltaPercent)
}

func TestCostComparison_Initialization(t *testing.T) {
	// Test initialization of CostComparison
	comparison := CostComparison{
		FirstMonthStart:  "2024-01-01",
		SecondMonthStart: "2024-02-01",
		IsProjection:     true,
		Multiplier:       1.0,
		ServiceCosts: []ServiceCosts{
			{
				ServiceName:  "Amazon EC2",
				Amount:       100.00,
				SecondAmount: 120.00,
				Delta:        20.00,
				DeltaPercent: 20.00,
			},
		},
		TotalAmount:       100.00,
		TotalSecondAmount: 120.00,
		TotalDelta:        20.00,
		TotalDeltaPercent: 10.00,
	}

	// Assertions
	assert.Equal(t, "2024-01-01", comparison.FirstMonthStart)
	assert.Equal(t, "2024-02-01", comparison.SecondMonthStart)
	assert.True(t, comparison.IsProjection)
	assert.Equal(t, 1.0, comparison.Multiplier)
	assert.Len(t, comparison.ServiceCosts, 1)
	assert.Equal(t, "Amazon EC2", comparison.ServiceCosts[0].ServiceName)
	assert.Equal(t, 100.00, comparison.TotalAmount)
	assert.Equal(t, 120.00, comparison.TotalSecondAmount)
	assert.Equal(t, 20.00, comparison.TotalDelta)
	assert.Equal(t, 10.00, comparison.TotalDeltaPercent)
}

func TestCostComparison_EmptyInitialization(t *testing.T) {
	// Test empty initialization of CostComparison
	comparison := CostComparison{}

	// Assertions
	assert.Empty(t, comparison.FirstMonthStart)
	assert.Empty(t, comparison.SecondMonthStart)
	assert.False(t, comparison.IsProjection)
	assert.Equal(t, 0.0, comparison.Multiplier)
	assert.Empty(t, comparison.ServiceCosts)
	assert.Equal(t, 0.0, comparison.TotalAmount)
	assert.Equal(t, 0.0, comparison.TotalSecondAmount)
	assert.Equal(t, 0.0, comparison.TotalDelta)
	assert.Equal(t, 0.0, comparison.TotalDeltaPercent)
}
