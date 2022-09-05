package executor

import (
	"fmt"
	"gotest.tools/assert"
	"sync"
	"testing"
	"time"
)

func TestThreadSafetyTaskPriorityQueue(t *testing.T) {
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
				task := NewTestSliceCollectingExecutable(value, collection, nil)

				// submit it
				queue.Push(task)
				wg.Done()
			},
			iterations: 100,
			processes:  3,
			collection: make([]int, 0),
		},
		{
			name: "priority queue push and pop is thread safe",
			testedFunction: func(queue TaskQueue, collection *[]int, value int, wg *sync.WaitGroup) {
				// create task
				task := NewTestSliceCollectingExecutable(value, collection, nil)

				// submit it
				queue.Push(task)

				// pop the task (it might not be the task we've enqueued)
				// since this operation is blocking we need the task to not be nil
				popped := queue.Pop()
				assert.Check(t, popped != nil)
				wg.Done()
			},
			iterations: 100,
			processes:  3,
			collection: make([]int, 0),
		},
		{
			name: "priority queue push, pop and contains is thread safe",
			testedFunction: func(queue TaskQueue, collection *[]int, value int, wg *sync.WaitGroup) {
				task := NewTestSliceCollectingExecutable(value, collection, nil)

				queue.Push(task)

				popped := queue.Pop()
				assert.Check(t, popped != nil)

				// ignore result
				queue.List()
				wg.Done()
			},
			iterations: 100,
			processes:  3,
			collection: make([]int, 0),
		},
		{
			name: "priority queue push, pop, contains, get by id is thread safe",
			testedFunction: func(queue TaskQueue, collection *[]int, value int, wg *sync.WaitGroup) {
				task := NewTestSliceCollectingExecutable(value, collection, nil)

				someId := queue.Push(task)

				queue.Get(someId)

				// ignore result
				queue.List()

				popped := queue.Pop()
				assert.Check(t, popped != nil)

				wg.Done()
			},
			iterations: 100,
			processes:  3,
			collection: make([]int, 0),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Println(fmt.Sprintf("Start: %v", test.name))

			queue := NewTaskQueue()
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
