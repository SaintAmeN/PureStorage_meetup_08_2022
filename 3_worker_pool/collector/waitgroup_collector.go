package collector

import (
	"AwesomePresentation/3_worker_pool/model"
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

type waitGroupCollector struct {
	nodes []model.ProcessingNode
}

func NewWaitGroupCollector(nodes []model.ProcessingNode) model.Collector {
	return &waitGroupCollector{
		nodes: nodes,
	}
}

func (c *waitGroupCollector) CollectResultsForValue(value float64) model.CollectionResult {
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

	// Collection (slice) to which we're going to collect results
	collectedResults := model.CollectionResult{}
	numberOfFailed := 0
	numberOfSuccessful := 0
	lock := &sync.Mutex{}
	waitGroup := &sync.WaitGroup{}

	// increase the wait group to expect given amount of results
	waitGroup.Add(len(c.nodes))

	// Iterate over all devices
	for _, node := range c.nodes {

		// Process each node in goroutine
		go func(processingNode model.ProcessingNode,
			sharedLock *sync.Mutex,
			sharedWaitGroup *sync.WaitGroup,
			sharedNumberOfSuccessful *int,
			sharedNumberOfFailed *int) {
			// Covers special case if runtime.Goexit() is called
			var err error = fmt.Errorf("goroutine exited before collection could finish")

			var output *float64
			defer (func() {
				defer sharedWaitGroup.Done()

				// If it panics in goroutine, we need to return error to channel
				if panic := recover(); panic != nil {
					fmt.Println("panicked in goroutine: %v \n\n %v", panic, string(debug.Stack()))
					err = fmt.Errorf("PANIC: %v", panic)
				}

				calculationOutputResult := model.CalculationOutput{
					Result: output,
					Error:  err,
				}

				sharedLock.Lock()
				defer sharedLock.Unlock()
				collectedResults = append(collectedResults, calculationOutputResult)
				if err != nil {
					*sharedNumberOfFailed++
				} else {
					*sharedNumberOfSuccessful++
				}
			})()

			output, err = processingNode.Calculate(ctxWithTimeout, model.CalculationInput{InputValue: value})
		}(node, lock, waitGroup, &numberOfSuccessful, &numberOfFailed)
	}

	waitGroup.Wait()
	fmt.Println(fmt.Sprintf("Collected results: %v, Collected errors: %v", numberOfSuccessful, numberOfFailed))

	return collectedResults
}
