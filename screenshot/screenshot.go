package screenshot

import (
	"github.com/danielhoward-me/chaos-backend/screenshot/utils"
	"github.com/danielhoward-me/chaos-backend/screenshot/worker"

	_ "embed"
)

type Request struct {
	Data string `json:"data"`
}

//go:embed placeholder.jpg
var PlaceholderImage []byte

func Queue(data string) int {
	hash := utils.Hash(data)

	waitTime := 0
	if !utils.Exists(hash) {
		waitTime = worker.GetEstimatedWaitTime()
		worker.AddJob(worker.Job{
			Data: data,
			Hash: hash,
		})
	}

	return waitTime
}

func GetEstimatedWaitTime() int {
	return worker.GetEstimatedWaitTime()
}
