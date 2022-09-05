package main

import (
	"bufio"
	"fmt"
	"os"
)

type numberSender struct {
	Channel chan string
}

func NewNumberSender(sharedChannel chan string) numberSender { // Publisher
	return numberSender{
		Channel: sharedChannel,
	}
}

func (s *numberSender) send(value string) {
	fmt.Printf("\nSender: Sending '%v' to channel\n", value)
	s.Channel <- value
}

type numberReceiver struct {
	Channel chan string
}

func NewNumberReceiver(sharedChannel chan string) numberReceiver { // Subscriber
	return numberReceiver{
		Channel: sharedChannel,
	}
}

func (s *numberReceiver) receive() {
	fmt.Printf("\nReceiver: Starting receiver channel\n")
	for {
		value := <-s.Channel

		fmt.Printf("\nReceiver: Received on channel: '%v'\n", value)
	}
}

func main() {
	sharedChannel := make(chan string)
	receiver := NewNumberReceiver(sharedChannel)

	// asynchronously start running the method in separate goroutine (thread)
	go receiver.receive()

	// This will allow us to enter numbers until we write -1
	reader := bufio.NewReader(os.Stdin)

	// As You can see sender CAN BE created much later...
	sender := NewNumberSender(sharedChannel)

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
		}
		sender.send(trimmedValue)
	}

	fmt.Print("Thank You and goodbye!") // Graceful finish :))
}
