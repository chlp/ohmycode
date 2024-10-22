package main

import (
	"ohmycode_api/config"
	"ohmycode_api/internal/api"
	"ohmycode_api/internal/store"
	"ohmycode_api/internal/worker"
)

func main() {
	apiConfig := config.LoadApiConf()
	apiStore := store.NewStore(apiConfig.DB)
	worker.NewWorker(apiStore).Run()
	api.NewService(apiStore).Run()
}
