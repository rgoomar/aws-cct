package models

// ServiceCosts represents the cost data for a service
type ServiceCosts struct {
	ServiceName  string
	Amount       float64
	SecondAmount float64
	Delta        float64
	DeltaPercent float64
}

// CostComparison represents the comparison between two time periods
type CostComparison struct {
	FirstMonthStart   string
	SecondMonthStart  string
	IsProjection      bool
	Multiplier        float64
	ServiceCosts      []ServiceCosts
	TotalAmount       float64
	TotalSecondAmount float64
	TotalDelta        float64
	TotalDeltaPercent float64
}
