package worker

import (
	"ohmycode_api/internal/store"
	"ohmycode_api/pkg/util"
	"time"
)

type Worker struct {
	fileStore *store.FileStore
}

func NewWorker(store *store.FileStore) *Worker {
	return &Worker{
		fileStore: store,
	}
}

const timeToSleepBetweenRuns = 100 * time.Millisecond

func (w *Worker) Run() {
	util.Log(nil, "Worker started")
	go func() {
		for {
			files := w.fileStore.GetAllFiles()
			for _, file := range files {
				file.CleanupUsers()
				file.CleanupWriter()
				file.CleanupWaitingForResult()

				// send insert and update into db

				if file.IsUnused() {
					w.fileStore.DeleteFile(file.ID)
				}
			}
			time.Sleep(timeToSleepBetweenRuns)
		}
	}()
}
