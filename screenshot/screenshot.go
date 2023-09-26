package screenshot

import (
	"github.com/danielhoward-me/chaos-backend/screenshot/utils"
	"github.com/danielhoward-me/chaos-backend/screenshot/worker"
)

type Request struct {
	Data string `json:"data"`
}

func Queue(data string) {
	hash := utils.Hash(data)

	if utils.Exists(hash) {
		return
	}

	worker.AddJob(worker.Job{
		Data: data,
		Hash: hash,
	})
}
