package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
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

type periodicNumberSender struct {
	Channel      chan string
	repeatPeriod time.Duration
}

func NewPeriodicNumberSender(sharedChannel chan string, repeatPeriod time.Duration) periodicNumberSender { // Publisher
	return periodicNumberSender{
		Channel:      sharedChannel,
		repeatPeriod: repeatPeriod,
	}
}

func (s *periodicNumberSender) send() {
	fmt.Printf("\nPeriodic sender: Starting 2_channels_with_periodic_thread sender with period %v", s.repeatPeriod)
	tick := time.NewTicker(s.repeatPeriod * time.Second)
	counter := 0
	for {
		select {
		case <-tick.C:
			fmt.Printf("\nPeriodic sender: Clock ticked, sending %v", counter)
			s.Channel <- string(counter)
			counter++
		}
	}
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
	fmt.Printf("\nReceiver: Starting receiver channel")
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

	// As You can see sender CAN BE created much later...
	periodicSender := NewPeriodicNumberSender(sharedChannel, 3)
	go periodicSender.send()

	fmt.Printf("\nMain: Starting 1_channels loop\n")
	for {
		fmt.Print("\nMain:Provide text, or 'quit' to finish:")
		value, err := reader.ReadString('\n')
		if err != nil {
			fmt.Print("Main: Something went really wrong...")
			return
		}

		// Let's remove last character (which is always end line '\n')
		sender.send(value[:len(value)-2])
	}

	fmt.Print("\n\nThank You and goodbye!") // Graceful finish :))
}
