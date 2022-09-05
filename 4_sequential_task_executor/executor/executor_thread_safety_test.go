package executor

import (
	"fmt"
	"gotest.tools/assert"
	"sync"
	"testing"
	"time"
)

func TestThreadSafetyExecutor(t *testing.T) {
	tests := []struct {
		name           string
		testedFunction func(queue TaskQueue, collection *[]int, value int, wg *sync.WaitGroup)
		iterations     int
		processes      int
		collection     []int
	}{
		{
			name: "priority queue push is thread safe",
			testedFunction: func(queue TaskQueue, collection *[]int, value int, wg *sync.WaitGroup) {
				// create task
				task := NewTestSliceCollectingExecutable(value, collection, wg)

				// submit it
				queue.Push(task)
			},
			iterations: 100,
			processes:  3,
			collection: make([]int, 0),
		},
		{
			name: "priority queue push and list is thread safe",
			testedFunction: func(queue TaskQueue, collection *[]int, value int, wg *sync.WaitGroup) {
				// create task
				task := NewTestSliceCollectingExecutable(value, collection, wg)

				// submit it
				queue.Push(task)

				queue.List()
			},
			iterations: 100,
			processes:  3,
			collection: make([]int, 0),
		},
		{
			name: "priority queue push, contains, get by id is thread safe",
			testedFunction: func(queue TaskQueue, collection *[]int, value int, wg *sync.WaitGroup) {
				task := NewTestSliceCollectingExecutable(value, collection, wg)

				// ignore result
				queue.List()
				someId := queue.Push(task)

				queue.Get(someId)

				// ignore result
				queue.List()
			},
			iterations: 100,
			processes:  3,
			collection: make([]int, 0),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Println(fmt.Sprintf("Start: %v", test.name))

			queue, _ := NewChannelQueueExecutor()
			//queue, _ := NewLockingQueueExecutor()

			setOfOperations := make([]ExecutableTestOperation, 0, test.iterations)
			wg := sync.WaitGroup{}
			wg.Add(test.iterations)

			for i := 0; i < test.iterations; i++ {
				copyOfFunction := test.testedFunction
				copyOfIndex := i

				operation := func() error {
					copyOfFunction(queue, &test.collection, copyOfIndex, &wg)
					return nil
				}

				setOfOperations = append(setOfOperations, operation)
			}

			fmt.Println(fmt.Sprintf("Testing number of operations: %v", len(setOfOperations)))
			_ = ParallelOperationsExecutor(t, test.processes, setOfOperations)

			result := WaitGroupWithTimeout(&wg, 5*time.Second)
			assert.Check(t, result)
		})
	}
}
