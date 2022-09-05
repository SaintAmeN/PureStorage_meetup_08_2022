package main

import (
	"AwesomePresentation/4_sequential_task_executor/executor"
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	// This will allow us to enter numbers until we write -1
	reader := bufio.NewReader(os.Stdin)

	// Create the queue
	queue, _ := executor.NewLockingQueueExecutor()

	fmt.Printf("\nMain: Starting 1_channels loop\n")
	for {
		fmt.Print("\nMain:Provide text, or 'quit' to finish: \n")
		value, err := reader.ReadString('\n')
		if err != nil {
			fmt.Print("Main: Something went really wrong...")
			return
		}

		// Let's remove last character (which is always end line '\n')
		trimmedValue := value[:len(value)-2]
		if trimmedValue == "quit" {
			break
		} else if trimmedValue == "count" {
			queue.Push(executor.NewExecutableCounterWithSleep(10, time.Duration(500)*time.Millisecond))
		} else if trimmedValue == "child" {
			queue.Push(executor.NewExecutableAnnoyingKid())
		} else if trimmedValue == "quickie" {
			queue.Push(executor.NewExecutableQuickie())
		}
	}

	fmt.Print("Thank You and goodbye!") // Graceful finish :))
}
