package display

import (
	"testing"

	"github.com/rgoomar/aws-cct/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestTableDisplay_Render(t *testing.T) {
	// Create test data
	comparison := &models.CostComparison{
		FirstMonthStart:  "2024-01-01",
		SecondMonthStart: "2024-02-01",
		IsProjection:     false,
		Multiplier:       1.0,
		ServiceCosts: []models.ServiceCosts{
			{
				ServiceName:  "Amazon EC2",
				Amount:       100.00,
				SecondAmount: 120.00,
				Delta:        20.00,
				DeltaPercent: 20.00,
			},
			{
				ServiceName:  "Amazon S3",
				Amount:       50.00,
				SecondAmount: 45.00,
				Delta:        -5.00,
				DeltaPercent: -10.00,
			},
		},
		TotalAmount:       150.00,
		TotalSecondAmount: 165.00,
		TotalDelta:        15.00,
		TotalDeltaPercent: 10.00,
	}

	// Test table output
	tableDisplay := NewTableDisplay()
	output := tableDisplay.Render(comparison, "table")

	// Basic assertions
	assert.Contains(t, output, "Amazon EC2")
	assert.Contains(t, output, "Amazon S3")
	assert.Contains(t, output, "TOTAL")
	assert.Contains(t, output, "$100.00")
	assert.Contains(t, output, "$120.00")
	assert.Contains(t, output, "$20.00")
	assert.Contains(t, output, "20.0%")

	// Test CSV output
	csvOutput := tableDisplay.Render(comparison, "csv")
	assert.Contains(t, csvOutput, "Service,2024-01-01,2024-02-01,Delta,Delta Percent")
	assert.Contains(t, csvOutput, "Amazon EC2,$100.00,$120.00,$20.00,20.0%")
	assert.Contains(t, csvOutput, "TOTAL,$150.00,$165.00,$15.00,10.0%")
}

func TestTableDisplay_RenderWithProjection(t *testing.T) {
	// Create test data with projection
	comparison := &models.CostComparison{
		FirstMonthStart:  "2024-01-01",
		SecondMonthStart: "2024-02-01",
		IsProjection:     true,
		Multiplier:       1.0,
		ServiceCosts: []models.ServiceCosts{
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
		TotalDeltaPercent: 20.00,
	}

	// Test table output with projection
	tableDisplay := NewTableDisplay()
	output := tableDisplay.Render(comparison, "table")

	// Assertions
	assert.Contains(t, output, "(PROJECTION)")
	assert.Contains(t, output, "Amazon EC2")
	assert.Contains(t, output, "$100.00")
	assert.Contains(t, output, "$120.00")
}
