package main

import (
	"context"
	"ohmycode_api/config"
	"ohmycode_api/internal/api"
	"ohmycode_api/internal/store"
	"ohmycode_api/internal/worker"
)

func main() {
	appCtx := context.Background()

	apiConfig := config.LoadApiConf()

	fileStore := store.NewFileStore(apiConfig.DB)
	runnerStore := store.NewRunnerStore()
	taskStore := store.NewTaskStore()

	worker.NewWorker(appCtx, fileStore).Run()

	api.NewService(apiConfig.HttpPort, fileStore, runnerStore, taskStore).Run()
}
