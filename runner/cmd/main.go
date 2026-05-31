package main

import (
	"context"
	"fmt"
	"ohmycode_runner/config"
	"ohmycode_runner/internal/api"
	"ohmycode_runner/internal/worker"
	"ohmycode_runner/pkg/util"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	appCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	apiConfig := config.LoadRunnerConf()
	util.Log(appCtx, fmt.Sprintf("OhMyCode.Runner app started with id: %s", apiConfig.RunnerId))
	apiClient := api.NewApiClient(appCtx, apiConfig.RunnerId, apiConfig.IsPublic, apiConfig.ApiUrl)

	w := worker.NewWorker(appCtx, apiConfig.RunnerId, apiClient, apiConfig.Languages)
	w.Run()

	<-appCtx.Done()
	util.Log(appCtx, "Shutting down: waiting for worker goroutines...")
	w.WaitDone()
	util.Log(appCtx, "Shutdown complete")
}
