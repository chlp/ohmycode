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
	timeToSleepBetweenCleanups          = 100 * time.Millisecond
	timeToSleepBetweenPersists          = 30 * time.Second
	timeToSleepBetweenSetIsRunnerOnline = 500 * time.Millisecond
)

func (w *Worker) Run() {
	util.Log("Worker started")
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
	go func() {
		for {
			select {
			case <-w.appCtx.Done():
				return
			default:
				w.filesSetIsRunnerOnline()
				time.Sleep(timeToSleepBetweenSetIsRunnerOnline)
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
