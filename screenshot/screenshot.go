package screenshot

import (
	"github.com/danielhoward-me/chaos-backend/screenshot/status"
	"github.com/danielhoward-me/chaos-backend/screenshot/utils"
	"github.com/danielhoward-me/chaos-backend/screenshot/worker"
)

type Request struct {
	Data     string `json:"data"`
	ForceNew bool   `json:"forceNew"`
}

func Queue(data string, forceNew bool) {
	hash := utils.Hash(data)

	if !forceNew && (utils.Exists(hash) || worker.InQueue(hash)) {
		return
	}

	worker.AddJob(worker.Job{
		Data: data,
		Hash: hash,
	})
}

func GetStatus(hash string) status.Status {
	if utils.Exists(hash) {
		return status.Generated
	} else if worker.HasFailed(hash) {
		return status.Failed
	} else if worker.IsProcessing(hash) {
		return status.Generating
	} else if worker.InQueue(hash) {
		return status.InQueue
	} else {
		return status.NotInQueue
	}
}
