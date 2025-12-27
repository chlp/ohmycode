package worker

import (
	"context"
	"ohmycode_api/internal/store"
	"ohmycode_api/pkg/util"
	"time"
)

type Worker struct {
	fileStore   *store.FileStore
	runnerStore *store.RunnerStore
	appCtx      context.Context
}

func NewWorker(appCtx context.Context, fileStore *store.FileStore, runnerStore *store.RunnerStore) *Worker {
	return &Worker{
		fileStore:   fileStore,
		runnerStore: runnerStore,
		appCtx:      appCtx,
	}
}

const (
	timeToSleepBetweenCleanups          = 1 * time.Second
	timeToSleepBetweenPersists          = 30 * time.Second
	timeToSleepBetweenSetIsRunnerOnline = 1 * time.Second
)

func (w *Worker) Run() {
	util.Log("Worker started")
	go func() {
		ticker := time.NewTicker(timeToSleepBetweenCleanups)
		defer ticker.Stop()
		for {
			select {
			case <-w.appCtx.Done():
				return
			case <-ticker.C:
				w.filesCleanUp()
			}
		}
	}()
	go func() {
		ticker := time.NewTicker(timeToSleepBetweenPersists)
		defer ticker.Stop()
		for {
			select {
			case <-w.appCtx.Done():
				return
			case <-ticker.C:
				w.filesPersisting()
			}
		}
	}()
	go func() {
		ticker := time.NewTicker(timeToSleepBetweenSetIsRunnerOnline)
		defer ticker.Stop()
		for {
			select {
			case <-w.appCtx.Done():
				return
			case <-ticker.C:
				w.filesSetIsRunnerOnline()
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
		if !file.Persisted {
			continue
		}
		if !file.UpdatedAt.After(file.PersistedAt) {
			continue
		}
		_ = w.fileStore.PersistFile(file)
	}
}

func (w *Worker) filesSetIsRunnerOnline() {
	files := w.fileStore.GetAllFiles()
	for _, file := range files {
		if file.UsePublicRunner {
			file.IsRunnerOnline = w.runnerStore.IsOnline(true, "")
		} // todo: implement for !UsePublicRunner
	}
}
