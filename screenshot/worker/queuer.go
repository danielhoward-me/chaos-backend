package worker

import (
	"fmt"
)

var chaosDevPort int
var jobs []Job
var workerRunning = false

type Job struct {
	Data string
	Hash string
}

func Init(givenChaosDevPort int) {
	chaosDevPort = givenChaosDevPort
}

func AddJob(job Job) {
	jobs = append(jobs, job)

	fmt.Printf("Queuing screenshot job for %s\n", job.Hash)

	if !workerRunning {
		go run()
	}
}

func GetEstimatedWaitTime() int {
	return len(jobs) * 10000
}
