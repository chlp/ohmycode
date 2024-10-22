package worker

import (
	"ohmycode_api/internal/store"
	"ohmycode_api/pkg/util"
	"time"
)

type Worker struct {
	store *store.Store
}

func NewWorker(store *store.Store) *Worker {
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
		}
		time.Sleep(timeToSleepBetweenRuns)
	}()
}
