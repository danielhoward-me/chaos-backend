package screenshot

import (
	"github.com/danielhoward-me/chaos-backend/screenshot/worker"

	_ "embed"
)

type Request struct {
	Data string `json:"data"`
}

//go:embed placeholder.jpg
var PlaceholderImage []byte

func QueueGeneration(data string) {
	worker.AddJob(data)
}

func GetEstimatedWaitTime() int {
	return worker.GetEstimatedWaitTime()
}
