package worker

import (
	"github.com/danielhoward-me/chaos-backend/screenshot/utils"

	"fmt"
)

var chaosDevPort int
var jobs []Job
var workerRunning = false

type Job struct {
	data string
	hash string
}

func Init(givenChaosDevPort int) {
	chaosDevPort = givenChaosDevPort
}

func AddJob(data string) {
	job := Job{
		data: data,
		hash: utils.Hash(data),
	}
	jobs = append(jobs, job)

	fmt.Printf("Queuing screenshot job for %s\n", job.hash)

	if !workerRunning {
		go run()
	}
}

func GetEstimatedWaitTime() int {
	return len(jobs) * 10000
}
