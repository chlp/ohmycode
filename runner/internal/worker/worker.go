package worker

import (
	"context"
	"ohmycode_runner/internal/api"
	"ohmycode_runner/pkg/util"
	"sync"
	"time"
)

type Worker struct {
	appCtx    context.Context
	runnerId  string
	apiClient *api.Client
	languages []string
	wg        sync.WaitGroup
}

const dataFolderPath = "data"

func NewWorker(appCtx context.Context, runnerId string, apiClient *api.Client, languages []string) *Worker {
	return &Worker{
		appCtx:    appCtx,
		runnerId:  runnerId,
		apiClient: apiClient,
		languages: languages,
	}
}

const intervalBetweenTasksReceive = 100 * time.Millisecond
const intervalBetweenResultsSend = 100 * time.Millisecond

func (w *Worker) Run() {
	util.Log(nil, "Worker started")
	taskDistributor := NewTaskDistributor(w.apiClient, w.runnerId, w.languages)
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		delay := intervalBetweenTasksReceive
		timer := time.NewTimer(delay)
		defer timer.Stop()
		for {
			select {
			case <-w.appCtx.Done():
				return
			case <-timer.C:
				processed, err := taskDistributor.Process()
				if err != nil || processed == 0 {
					delay = intervalBetweenTasksReceive * 10
				} else {
					delay = intervalBetweenTasksReceive
				}
				timer.Reset(delay)
			}
		}
	}()
	for _, language := range w.languages {
		resultProcessor := NewResultProcessor(w.apiClient, w.runnerId, language)
		w.wg.Add(1)
		go func(language string) {
			defer w.wg.Done()
			delay := intervalBetweenResultsSend
			timer := time.NewTimer(delay)
			defer timer.Stop()
			for {
				select {
				case <-w.appCtx.Done():
					return
				case <-timer.C:
					processed, err := resultProcessor.Process()
					if err != nil || processed == 0 {
						delay = intervalBetweenResultsSend * 10
					} else {
						delay = intervalBetweenResultsSend
					}
					timer.Reset(delay)
				}
			}
		}(language)
	}
}

func (w *Worker) WaitDone() {
	w.wg.Wait()
}
