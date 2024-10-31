package worker

import (
	"context"
	"ohmycode_runner/internal/api"
	"ohmycode_runner/pkg/util"
	"time"
)

type Worker struct {
	appCtx    context.Context
	runnerId  string
	apiClient *api.Client
	languages []string
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
	go func() {
		for {
			select {
			case <-w.appCtx.Done():
				return
			default:
				taskDistributor.Process()
				time.Sleep(intervalBetweenTasksReceive)
			}
		}
	}()
	for _, language := range w.languages {
		resultProcessor := NewResultProcessor(w.apiClient, w.runnerId, language)
		go func(language string) {
			for {
				select {
				case <-w.appCtx.Done():
					return
				default:
					resultProcessor.Process()
					time.Sleep(intervalBetweenResultsSend)
				}
			}
		}(language)
	}
}
