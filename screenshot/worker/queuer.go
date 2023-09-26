package worker

import (
	"fmt"
)

var chaosDevPort int

var jobs []Job
var workerRunning = false
var currentlyProcessing string

var failedScreenshots = map[string]bool{}

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

func InQueue(hash string) bool {
	for _, job := range jobs {
		if hash == job.Hash {
			return true
		}
	}

	return false
}

func IsProcessing(hash string) bool {
	if !workerRunning {
		return false
	}

	return currentlyProcessing == hash
}

func HasFailed(hash string) bool {
	failed := failedScreenshots[hash]

	// Set the screenshot to allow another request to be made
	// since a respecting client will stop making requests once
	// it it told it has failed
	delete(failedScreenshots, hash)

	return failed
}
