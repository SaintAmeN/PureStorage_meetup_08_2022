package node

import (
	"AwesomePresentation/3_worker_pool/model"
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// The purpose of this struct/object is to calculate value and output result divided by the `factor`
type divideProcessingNode struct {
	factor float64
}

func NewDivideProcessingNode(factor float64) model.ProcessingNode {
	return &divideProcessingNode{
		factor: factor,
	}
}

func (p *divideProcessingNode) Calculate(ctx context.Context, input model.CalculationInput) (*float64, error) {
	randomTime := (rand.Float64() * 10) + 1
	randomDuration := time.Duration(math.Abs(randomTime)) * time.Second

	fmt.Println(fmt.Sprintf("Multiply Processing Node: Picked random time : %v, duration: %v", randomTime, randomDuration))
	tick := time.NewTicker(randomDuration)

	select {
	case <-tick.C:
		result := input.InputValue / p.factor
		fmt.Println(fmt.Sprintf("Divide Processing Node: Returning : %v", result))

		return &result, nil
	case <-ctx.Done():
		msg := "Divide Processing Node: Timed out"
		fmt.Println(msg)
		return nil, fmt.Errorf(msg)
	}
}
