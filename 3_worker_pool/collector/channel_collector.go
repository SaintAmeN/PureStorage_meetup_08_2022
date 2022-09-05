package collector

import (
	"AwesomePresentation/3_worker_pool/model"
	"context"
	"fmt"
	"runtime/debug"
	"time"
)

type channelCollector struct {
	nodes []model.ProcessingNode
}

func NewChannelCollector(nodes []model.ProcessingNode) model.Collector {
	return &channelCollector{
		nodes: nodes,
	}
}

func (c *channelCollector) CollectResultsForValue(value float64) model.CollectionResult {
	// This is creating new context
	ctx := context.TODO()

	// Create context which will time out after configured amount of time
	ctxWithTimeout, cancelFunc := context.WithTimeout(ctx, time.Duration(5)*time.Second)
	defer cancelFunc()
	defer func() {
		if ctxErr := ctxWithTimeout.Err(); ctxErr != nil {
			fmt.Println(fmt.Errorf("Collect: Unexpected context error (adjust timeout accordingly): %v", ctxErr))
		}
	}()

	// Channel to receive commits from devices
	resultsChannel := make(chan model.CalculationOutput)

	// Iterate over all devices
	for _, node := range c.nodes {

		// Process each node in goroutine
		go func(processingNode model.ProcessingNode) {
			// Covers special case if runtime.Goexit() is called
			var err error = fmt.Errorf("goroutine exited before collection could finish")

			var output *float64
			defer (func() {
				// If it panics in goroutine, we need to return error to channel
				if panic := recover(); panic != nil {
					fmt.Println("panicked in goroutine: %v \n\n %v", panic, string(debug.Stack()))
					err = fmt.Errorf("PANIC: %v", panic)
				}

				resultsChannel <- model.CalculationOutput{
					Result: output,
					Error:  err,
				}
			})()

			output, err = processingNode.Calculate(ctxWithTimeout, model.CalculationInput{InputValue: value})
		}(node)
	}

	collectedResults := model.CollectionResult{}
	numberOfFailed := 0
	numberOfSuccessful := 0
	for range c.nodes {
		result := <-resultsChannel
		collectedResults = append(collectedResults, result)
		if result.Error != nil {
			numberOfFailed++
			fmt.Println(fmt.Sprintf("Result not gathered, error: %v", result.Error))
		} else {
			numberOfSuccessful++
			fmt.Println(fmt.Sprintf("Result: %v", *result.Result))
		}
	}
	fmt.Println(fmt.Sprintf("Collected results: %v, Collected errors: %v", numberOfSuccessful, numberOfFailed))

	return collectedResults
}
