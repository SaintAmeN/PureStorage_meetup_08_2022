package executor

import (
	"fmt"
	"github.com/google/uuid"
	"runtime/debug"
)

type taskQueue struct {
	tasks []Task

	requestChannel  chan queueRequest
	responseChannel chan queueResponse

	// We keep those separate because there are cases we're going to ignore the call
	executorRequestChannel  chan queueRequest
	executorResponseChannel chan queueResponse
}
type TaskQueue interface {
	// Fetches task from queue if there is one present
	Pop() *Task

	// Pushes new task to the queue, task is copied in the method. Returns task id
	Push(task Task) string

	// Lists all available tasks (creates copy of all tasks)
	List() []TaskInfo

	// Fetches task that supposedly exists in the queue by it's ID. Might return nil if task wasn't found
	Get(id string) *TaskInfo
}

func NewTaskQueue() TaskQueue {
	taskQueue := &taskQueue{
		tasks:                   make([]Task, 0),
		requestChannel:          make(chan queueRequest),
		responseChannel:         make(chan queueResponse),
		executorRequestChannel:  make(chan queueRequest),
		executorResponseChannel: make(chan queueResponse),
	}

	go taskQueue.RunQueue()
	return taskQueue
}

func (q *taskQueue) Pop() *Task {
	q.executorRequestChannel <- queueGetTaskRequest{}

	response := <-q.executorResponseChannel

	var result Task

	switch castedResponse := response.(type) {
	case queueGetTaskResponse:
		result = castedResponse.task
	default:
		fmt.Println(fmt.Errorf("failed to pop task, incorrect type"))
	}

	return &result
}

func (q *taskQueue) Push(task Task) string {
	q.requestChannel <- queueEnqueueTaskRequest{task: task}

	response := <-q.responseChannel

	var result string

	switch castedResponse := response.(type) {
	case queueEnqueueTaskResponse:
		result = castedResponse.taskId
	default:
		fmt.Println(fmt.Errorf("failed to push task, incorrect type"))
	}

	return result
}

func (q *taskQueue) List() []TaskInfo {
	q.requestChannel <- queueGetListOfTasksRequest{}

	response := <-q.responseChannel

	var result []TaskInfo

	switch castedResponse := response.(type) {
	case queueGetListOfTasksResponse:
		result = castedResponse.tasks
	default:
		fmt.Println(fmt.Errorf("failed to fetch list of tasks, incorrect type"))
	}

	return result
}

func (q *taskQueue) Get(id string) *TaskInfo {
	q.requestChannel <- queueGetTaskByIdRequest{taskId: id}

	response := <-q.responseChannel

	var result *TaskInfo

	switch castedResponse := response.(type) {
	case queueGetTaskByIdResponse:
		result = castedResponse.task
	default:
		fmt.Println(fmt.Errorf("failed to fetch task by id, incorrect type"))
	}

	return result
}

func (q *taskQueue) RunQueue() {
	defer (func() {
		if panic := recover(); panic != nil {
			fmt.Println(fmt.Errorf("queue receiver goroutine panicked: %v \n\n %v", panic, string(debug.Stack())))
		}
	})()

	for {
		fmt.Println(fmt.Errorf("queue iteration, current length is: %v awaiting requests", len(q.tasks)))
		if len(q.tasks) > 0 {
			// if there are elements enqueued we await both incoming and outgoing messages
			select {
			case request := <-q.executorRequestChannel:
				// if there is active listener on outgoing channel, pop top of the queue
				q.processExecutorRequest(request)
			case request := <-q.requestChannel:
				// accept incoming requests and process them
				q.processRequest(request)
			}
		} else {
			// if the queue is empty we await incoming messages
			request := <-q.requestChannel
			q.processRequest(request)
		}
	}
}

// Method used by the `receiver` goroutine. Processing `requestChannel` channel is done by
//	pulling requests from it and handling them correctly. Goroutine is handling three
//	types of requests: `enqueue`, `stop` (aka cancel), `containsTaskWithPriority`
func (q *taskQueue) processExecutorRequest(request queueRequest) {
	fmt.Println(fmt.Sprintf("queue received executor request: %v", request))

	switch req := request.(type) {
	case queueGetTaskRequest:
		q.processQueueGetTaskRequest()
	default:
		// we need to handle default not to be blocked
		fmt.Println(fmt.Errorf("queue received invalid/unknown executor request type: %v discarded", req))
	}
}

func (q *taskQueue) processRequest(request queueRequest) {
	fmt.Println(fmt.Sprintf("queue received request: %v", request))

	switch req := request.(type) {
	case queueGetTaskByIdRequest:
		q.processQueueGetTaskByIdRequest(req)
	case queueGetListOfTasksRequest:
		q.processQueueGetListOfTasksRequest(req)
	case queueEnqueueTaskRequest:
		q.processQueueEnqueueTaskRequest(req)
	default:
		// we need to handle default not to be blocked
		fmt.Println(fmt.Errorf("queue received invalid/unknown request type: %v discarded", req))
	}
}

func (q *taskQueue) processQueueGetTaskRequest() {
	fmt.Println(fmt.Sprintf("queue called pop"))
	if len(q.tasks) == 0 {
		fmt.Println(fmt.Errorf("queue called pop with empty queue"))
		return
	}

	firstTask := q.tasks[0]
	q.tasks = q.tasks[1:]

	q.executorResponseChannel <- queueGetTaskResponse{firstTask}
}

func (q *taskQueue) processQueueEnqueueTaskRequest(request queueEnqueueTaskRequest) {
	fmt.Println(fmt.Sprintf("queue called enqueue"))

	generatedTaskId, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(fmt.Errorf("failed to generate uuid for operation: %w", err))
	}

	taskIdString := generatedTaskId.String()
	copiedTask := request.task
	copiedTask.Id = taskIdString

	q.tasks = append(q.tasks, copiedTask)
	q.responseChannel <- queueEnqueueTaskResponse{taskId: taskIdString}
}

func (q *taskQueue) processQueueGetListOfTasksRequest(req queueGetListOfTasksRequest) {
	fmt.Println(fmt.Sprintf("queue called get list of tasks"))

	tasksCopy := make([]TaskInfo, 0)
	for _, task := range q.tasks {
		tasksCopy = append(tasksCopy, TaskInfo{
			Id: task.Id,
		})
	}

	q.responseChannel <- queueGetListOfTasksResponse{tasks: tasksCopy}
}

func (q *taskQueue) processQueueGetTaskByIdRequest(req queueGetTaskByIdRequest) {
	fmt.Println(fmt.Sprintf("queue called get task by id"))

	var returnedTask *TaskInfo = nil
	for _, taskInQueue := range q.tasks {
		if taskInQueue.Id == req.taskId {
			returnedTask = &TaskInfo{
				Id: taskInQueue.Id,
			}

			fmt.Println(fmt.Sprintf("get task by id - found"))
			break
		}
	}

	q.responseChannel <- queueGetTaskByIdResponse{task: returnedTask}
}

// Queue Requests
// queueRequest is an internal interface (used only within this class)
type queueRequest interface{}
type queueGetTaskRequest struct{}
type queueEnqueueTaskRequest struct{ task Task }
type queueGetListOfTasksRequest struct{}
type queueGetTaskByIdRequest struct{ taskId string }

// Queue Requests
// queueRequest is an internal interface (used only within this class)
type queueResponse interface{}
type queueGetTaskResponse struct{ task Task }
type queueEnqueueTaskResponse struct{ taskId string }
type queueGetListOfTasksResponse struct{ tasks []TaskInfo }
type queueGetTaskByIdResponse struct{ task *TaskInfo }
