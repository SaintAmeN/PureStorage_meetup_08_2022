package executor

import (
	"fmt"
	"gotest.tools/assert"
	"runtime"
	"sync"
	"testing"
	"time"
)

// It might happen that the parallel test executes with success couple of times, so we need to run it in a loop just for confirmation
const DefaultNumberOfTestIterations = 1

type ExecutableTestOperation func() error

// Wrapper for custom function. Parallel operations are being tested multiple times to make sure we're not creating flaky tests
func ParallelOperationsExecutor(t *testing.T, numberOfParallelRoutines int, testedOperations []ExecutableTestOperation) []error {
	return CustomParallelOperationsExecutor(t, numberOfParallelRoutines, testedOperations, DefaultNumberOfTestIterations)
}

func CustomParallelOperationsExecutor(t *testing.T, numberOfParallelRoutines int, testedOperations []ExecutableTestOperation, numberOfExecutions int) []error {
	runtime.GOMAXPROCS(numberOfParallelRoutines)

	var lastExecutionResult []error
	for i := 0; i < numberOfExecutions; i++ {
		// If execution in loop succeeded we're not really interested in results
		// If we are, then this loop will only return the latest ones
		lastExecutionResult = executeParallelOperations(t, numberOfParallelRoutines, testedOperations)
	}
	return lastExecutionResult
}

func executeParallelOperations(t *testing.T, numberOfParallelRoutines int, testedOperations []ExecutableTestOperation) []error {
	expectedTestedConcurrentOperationsCount := len(testedOperations)

	// This lock releases execution of 1_channels thread once all goroutines have been started
	// Because it takes time to fire goroutine (not much, but some) we need to make sure we wait for all of them to be alive
	// Once all of them are ready, 1_channels thread is notified and we continue execution
	startupLock := &sync.WaitGroup{}
	startupLock.Add(expectedTestedConcurrentOperationsCount)

	// Execution for all goroutines is continued when 1_channels thread is notified that all goroutines are ready to continue (see above)
	parallelReleaseLock := &sync.WaitGroup{}
	parallelReleaseLock.Add(1)

	// Main thread should wait for all goroutines to finish, once they are done 1_channels thread is notified and we can return results
	completionLock := &sync.WaitGroup{}
	completionLock.Add(expectedTestedConcurrentOperationsCount)

	// Mutex to secure critical section (at the end of each goroutine)
	mutex := &sync.Mutex{}

	counter := 0
	errorResults := make([]error, 0)
	for index, operation := range testedOperations {
		go func(currentIdx int, currentOperation ExecutableTestOperation, counterPointer *int) {
			// Regardless of the result, we want the 1_channels thread to be notified and finally release execution
			defer completionLock.Done()

			// Notify 1_channels thread You've started
			startupLock.Done()

			// Wait for 1_channels thread to notify goroutines that they can continue
			// scp.Log().Infof("starting routine: %v waiting for parallelReleaseLock...", currentIdx)
			parallelReleaseLock.Wait()

			// scp.Log().Infof("executing operation index: %v", currentIdx)
			err := currentOperation()
			fmt.Println(fmt.Sprintf("finished executing operation index: %v", currentIdx))

			// Error while executing operation
			if err != nil {
				fmt.Println(fmt.Sprintf("execution error: %v operation: %v", err, currentIdx))
			}

			// Utility itself should be thread safe
			mutex.Lock()
			// we increase our counter
			// because storing locks, we can be sure that only one thread/routine at the time will reach this instruction
			// if implementation was incorrect, we'd experience race condition and failures.
			*counterPointer++
			errorResults = append(errorResults, err)
			mutex.Unlock()

		}(index, operation, &counter)
	}

	// We wait for a second to make sure all routines had time to reach locking point
	startupLock.Wait()

	fmt.Println("Waiting for routines finished, releasing parallelReleaseLock...")
	// Releasing lock should fire all routines to continue
	parallelReleaseLock.Done()

	// We wait for all routines to complete
	completionLock.Wait()
	fmt.Println("All routines completed...")

	// This check actually verifies that while executing we've locked in the middle and/or we didn't have race condition
	assert.Equal(t, expectedTestedConcurrentOperationsCount, counter)
	assert.Equal(t, expectedTestedConcurrentOperationsCount, len(errorResults))

	return errorResults
}

type testSliceCollectingExecutable struct {
	value      int
	collection *[]int
	wg         *sync.WaitGroup
}

func (e *testSliceCollectingExecutable) Execute() error {
	if e.wg != nil {
		// we need to wait for each task to finish, this way we can be sure test can terminate once all executable tasks are done, not before
		defer e.wg.Done()
	}
	// If there is a race condition, this statement will detect it.
	*e.collection = append(*e.collection, e.value)
	return nil
}

func NewTestSliceCollectingExecutable(value int, channelledCollectSlice *[]int, waitGroup *sync.WaitGroup) Task {
	return Task{
		TaskExecutable: &testSliceCollectingExecutable{
			value:      value,
			collection: channelledCollectSlice,
			wg:         waitGroup,
		},
	}
}

func WaitGroupWithTimeout(wg *sync.WaitGroup, timeout time.Duration) (succeeded bool) {
	c := make(chan struct{})

	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		// (successful) normal completion without timeout
		return true
	case <-time.After(timeout):
		// (failed) completion with timeout
		return false
	}
}
