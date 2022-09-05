package main

import (
	"AwesomePresentation/3_worker_pool/collector"
	"AwesomePresentation/3_worker_pool/model"
	"AwesomePresentation/3_worker_pool/node"
	"fmt"
	"time"
)

func main() {
	nodes := []model.ProcessingNode{
		node.NewDivideProcessingNode(1),
		node.NewDivideProcessingNode(2),
		node.NewDivideProcessingNode(3),
		node.NewDivideProcessingNode(3),
		node.NewMultiplyProcessingNode(3),
		node.NewMultiplyProcessingNode(2),
		node.NewMultiplyProcessingNode(1),
	}
	//resultsCollector := collector.NewLockingCollector(nodes)
	//resultsCollector := collector.NewWaitGroupCollector(nodes)
	resultsCollector := collector.NewChannelCollector(nodes)

	resultsCollector.CollectResultsForValue(1)
	fmt.Print("\nFinished Round 1\n\n\n") // Finish, round 1
	time.Sleep(time.Duration(5) * time.Second)

	resultsCollector.CollectResultsForValue(2)
	fmt.Print("\nFinished Round 2\n\n\n") // Finish, round 2
	time.Sleep(time.Duration(5) * time.Second)

	resultsCollector.CollectResultsForValue(3)
	fmt.Print("\nFinished Round 3\n\n\n")   // Finish, round 3
	fmt.Print("\n\nThank You and goodbye!") // Graceful finish :))
}
