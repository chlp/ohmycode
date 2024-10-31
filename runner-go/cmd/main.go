package main

import (
	"context"
	"ohmycode_runner/config"
	"ohmycode_runner/internal/api"
	"ohmycode_runner/internal/worker"
	"ohmycode_runner/pkg/util"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	appCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	apiConfig := config.LoadRunnerConf()
	apiClient := api.NewApiClient(apiConfig.RunnerId, apiConfig.IsPublic, apiConfig.ApiUrl)

	worker.NewWorker(appCtx, apiConfig.RunnerId, apiClient, apiConfig.Languages).Run()

	<-appCtx.Done()
	util.Log(appCtx, "Application stopped")
	time.Sleep(2 * time.Second)
}
