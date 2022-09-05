package executor

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
)

type lockingTaskQueue struct {
	tasks []Task

	lock sync.Mutex
}

func NewLockingTaskQueue() TaskQueue {
	queue := &lockingTaskQueue{
		tasks: make([]Task, 0),
		lock:  sync.Mutex{},
	}

	return queue
}

func (q *lockingTaskQueue) Pop() *Task {
	q.lock.Lock()
	defer q.lock.Unlock()

	fmt.Println(fmt.Sprintf("queue called pop"))
	if len(q.tasks) == 0 {
		fmt.Println(fmt.Errorf("queue called pop with empty queue"))
		return nil
	}

	firstTask := q.tasks[0]
	q.tasks = q.tasks[1:]

	return &firstTask
}

func (q *lockingTaskQueue) Push(task Task) string {
	q.lock.Lock()
	defer q.lock.Unlock()
	fmt.Println(fmt.Sprintf("queue called enqueue"))

	generatedTaskId, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(fmt.Errorf("failed to generate uuid for operation: %w", err))
	}

	taskIdString := generatedTaskId.String()
	copiedTask := task
	copiedTask.Id = taskIdString

	q.tasks = append(q.tasks, copiedTask)

	return taskIdString
}

func (q *lockingTaskQueue) List() []TaskInfo {
	q.lock.Lock()
	defer q.lock.Unlock()

	fmt.Println(fmt.Sprintf("queue called get list of tasks"))

	tasksCopy := make([]TaskInfo, 0)
	for _, task := range q.tasks {
		tasksCopy = append(tasksCopy, TaskInfo{
			Id: task.Id,
		})
	}

	return tasksCopy
}

func (q *lockingTaskQueue) Get(id string) *TaskInfo {
	q.lock.Lock()
	defer q.lock.Unlock()
	fmt.Println(fmt.Sprintf("queue called get task by id"))

	var returnedTask *TaskInfo = nil
	for _, taskInQueue := range q.tasks {
		if taskInQueue.Id == id {
			returnedTask = &TaskInfo{
				Id: taskInQueue.Id,
			}

			fmt.Println(fmt.Sprintf("get task by id - found"))
			break
		}
	}

	return returnedTask
}
