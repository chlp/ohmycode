package worker

import (
	"ohmycode_api/internal/store"
	"ohmycode_api/pkg/util"
	"time"
)

type Worker struct {
	fileStore *store.FileStore
}

func NewWorker(fileStore *store.FileStore) *Worker {
	return &Worker{
		fileStore: fileStore,
	}
}

const timeToSleepBetweenCleanups = 100 * time.Millisecond
const timeToSleepBetweenPersists = 30 * time.Second

func (w *Worker) Run() {
	util.Log(nil, "Worker started")
	go func() {
		for {
			files := w.fileStore.GetAllFiles()
			for _, file := range files {
				file.CleanupUsers()
				file.CleanupWriter()
				file.CleanupWaitingForResult()

				if file.IsUnused() {
					w.fileStore.DeleteFile(file.ID)
				}
			}
			time.Sleep(timeToSleepBetweenCleanups)
		}
	}()
	go func() {
		for {
			files := w.fileStore.GetAllFiles()
			for _, file := range files {
				_ = w.fileStore.PersistFile(file)
			}
			time.Sleep(timeToSleepBetweenPersists)
		}
	}()
}
