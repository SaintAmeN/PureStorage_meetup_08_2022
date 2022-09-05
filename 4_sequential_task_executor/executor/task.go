package executor

import (
	"fmt"
	"math/rand"
	"time"
)

type TaskInfo struct {
	Id string
	// FIXME: could also have more data copied from Task
	//		We usually keep lots more information, such as:
	//			- the result
	//			- wait group (to be able to wait for the task to finish and block),
	//			- etc.
}

type Task struct {
	Id             string
	TaskExecutable Executable
}

type Executable interface {
	Execute() error
}

type ExecutableCounterWithSleep struct {
	CountLimit  int
	CountPeriod time.Duration
}

func NewExecutableCounterWithSleep(countLimit int, countPeriod time.Duration) Task {
	return Task{
		TaskExecutable: &ExecutableCounterWithSleep{
			CountLimit:  countLimit,
			CountPeriod: countPeriod,
		},
	}
}

func (e *ExecutableCounterWithSleep) Execute() error {
	fmt.Println(fmt.Sprintf("Counting starting"))
	for i := 0; i < e.CountLimit; i++ {
		time.Sleep(e.CountPeriod)
		fmt.Println(fmt.Sprintf("Counting: %v", i))
	}
	fmt.Println(fmt.Sprintf("Counting finished"))
	return nil
}

type ExecutableAnnoyingKid struct {
	Task
	RandomSentences []string
}

func (e *ExecutableAnnoyingKid) Execute() error {
	fmt.Println(fmt.Sprintf("Annoying kid stating"))
	for _, sentence := range e.RandomSentences {
		fmt.Println(sentence)

		// Sleep random amount of time
		time.Sleep(time.Duration(rand.Float64()*10.0) * time.Second)

		// Fail randomly with 10% chance
		if rand.Float64() < 0.1 {
			return fmt.Errorf("I am not sorry, I failed :) ")
		}
	}
	fmt.Println(fmt.Sprintf("Annoying kid finished"))
	return nil
}

func NewExecutableAnnoyingKid() Task {
	return Task{
		TaskExecutable: &ExecutableAnnoyingKid{
			RandomSentences: []string{
				"Are we there yet?",
				"Are we there yet?",
				"I need to go to the toilet.",
				"Where?",
				"Are we there yet?",
				"Are we there yet?",
				"Where?",
				"Where?",
				"Are we there yet?",
				"Are we there yet?",
				"Where?",
				"Why?",
				"I need to go to the toilet.",
				"Are we there yet?",
				"Are we there yet?",
				"What is that?",
				"What is this?",
				"Why?",
				"Where?",
			},
		},
	}
}

type ExecutableQuickie struct {
	Task
}

func (e *ExecutableQuickie) Execute() error {
	fmt.Println(fmt.Sprintf("Quickie stating"))

	fmt.Println(fmt.Sprintf(":))"))

	fmt.Println(fmt.Sprintf("Quickie finished"))
	return nil
}

func NewExecutableQuickie() Task {
	return Task{
		TaskExecutable: &ExecutableQuickie{},
	}
}
