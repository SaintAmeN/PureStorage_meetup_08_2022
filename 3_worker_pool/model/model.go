package model

import "context"

type CalculationInput struct {
	InputValue float64
}

type CalculationOutput struct {
	Result *float64
	Error  error
}

type CollectionResult []CalculationOutput

type ProcessingNode interface {
	Calculate(ctx context.Context, input CalculationInput) (*float64, error)
}

type Collector interface {
	CollectResultsForValue(value float64) CollectionResult
}
