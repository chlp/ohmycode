package worker

import (
	"ohmycode_api/internal/store"
	"ohmycode_api/pkg/util"
	"time"
)

type Worker struct {
	store *store.FileStore
}

func NewWorker(store *store.FileStore) *Worker {
	return &Worker{
		store: store,
	}
}

const timeToSleepBetweenRuns = 100 * time.Millisecond

func (w *Worker) Run() {
	util.Log(nil, "Worker started")
	go func() {
		files := w.store.GetAllFiles()
		for _, file := range files {
			file.CleanupUsers()
			file.CleanupWriter()
			file.CleanupWaitingForResult()

			// send insert and update into db

			// remove files from memory that not in usage anymore
		}
		time.Sleep(timeToSleepBetweenRuns)
	}()
}
