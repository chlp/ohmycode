package worker

import (
	"context"
	"ohmycode_api/internal/store"
	"ohmycode_api/pkg/util"
	"time"
)

type Worker struct {
	fileStore *store.FileStore
	appCtx    context.Context
}

func NewWorker(appCtx context.Context, fileStore *store.FileStore) *Worker {
	return &Worker{
		fileStore: fileStore,
		appCtx:    appCtx,
	}
}

const timeToSleepBetweenCleanups = 100 * time.Millisecond
const timeToSleepBetweenPersists = 30 * time.Second

func (w *Worker) Run() {
	util.Log(nil, "Worker started")
	go func() {
		for {
			select {
			case <-w.appCtx.Done():
				return
			default:
				w.filesCleanUp()
				time.Sleep(timeToSleepBetweenCleanups)
			}
		}
	}()
	go func() {
		for {
			select {
			case <-w.appCtx.Done():
				return
			default:
				w.filesPersisting()
				time.Sleep(timeToSleepBetweenPersists)
			}
		}
	}()
}

func (w *Worker) filesCleanUp() {
	files := w.fileStore.GetAllFiles()
	for _, file := range files {
		file.CleanupUsers()
		file.CleanupWriter()
		file.CleanupWaitingForResult()

		if file.IsUnused() {
			w.fileStore.DeleteFile(file.ID)
		}
	}
}

func (w *Worker) filesPersisting() {
	files := w.fileStore.GetAllFiles()
	for _, file := range files {
		if !file.UpdatedAt.After(file.PersistedAt) {
			continue
		}
		_ = w.fileStore.PersistFile(file)
	}
}
