package executor

import (
	"fmt"
	"runtime/debug"
	"time"
)

type Executor struct {
	queue TaskQueue
}

func NewChannelQueueExecutor() (TaskQueue, Executor) {
	taskQueue := NewTaskQueue()
	// At this point queue is already running

	executor := Executor{queue: taskQueue}
	go executor.runExecutor()
	// Starting executor

	return taskQueue, executor
}

func NewLockingQueueExecutor() (TaskQueue, Executor) {
	taskQueue := NewLockingTaskQueue()

	executor := Executor{queue: taskQueue}
	go executor.runExecutor()
	// Starting executor

	return taskQueue, executor
}

func (e *Executor) runExecutor() {
	var task *Task
	defer (func() {
		if panic := recover(); panic != nil {
			fmt.Println(fmt.Errorf("executor goroutine panicked: %v \n\n %v", panic, string(debug.Stack())))
		}
	})()

	for {
		task = e.queue.Pop()

		if task != nil {
			fmt.Println(fmt.Sprintf("found task: %v", task))

			// TODO: consider passing cancellable context
			err := task.TaskExecutable.Execute()
			if err != nil {
				fmt.Println(fmt.Errorf("finished execution of task with error: %v", err))
			}

			fmt.Println(fmt.Sprintf("finished execution of task: %v", task))
		} else {
			fmt.Println(fmt.Sprintf("Returned nil task, need to wait for task to appear"))
			time.Sleep(1 * time.Second)
		}
	}
}
