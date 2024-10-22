package main

import (
	"ohmycode_api/config"
	"ohmycode_api/internal/api"
	"ohmycode_api/internal/store"
	"ohmycode_api/internal/worker"
)

func main() {
	apiConfig := config.LoadApiConf()
	fileStore := store.NewFileStore(apiConfig.DB)
	runnerStore := store.NewRunnerStore()
	worker.NewWorker(fileStore).Run()
	api.NewService(fileStore, runnerStore).Run()
}
