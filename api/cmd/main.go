package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"ohmycode_api/config"
	"ohmycode_api/internal/api"
	"ohmycode_api/internal/model"
	"ohmycode_api/internal/store"
	"ohmycode_api/internal/worker"
	"syscall"
	"time"
)

func main() {
	appCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	apiConfig := config.LoadApiConf()
	model.SetContentMaxLength(apiConfig.ContentMaxLengthKb * 1024)

	versionStore := store.NewVersionStore(apiConfig.DB)
	fileStore := store.NewFileStore(apiConfig.DB, versionStore)
	runnerStore := store.NewRunnerStore()
	taskStore := store.NewTaskStore()

	worker.NewWorker(appCtx, fileStore, runnerStore).Run()

	svc := api.NewService(apiConfig.HttpPort, apiConfig.ServeClientFiles, apiConfig.UseDynamicFiles, apiConfig.WsAllowedOrigins, fileStore, runnerStore, taskStore, versionStore)
	if err := svc.Run(appCtx); err != nil && !errors.Is(err, context.Canceled) {
		// avoid log.Fatal to allow defer cleanup
		panic(err)
	}

	log.Println("Shutting down: flushing dirty files to MongoDB...")
	if err := fileStore.FlushAll(); err != nil {
		log.Println("FlushAll error:", err)
	}

	closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = fileStore.Close(closeCtx)
	_ = versionStore.Close(closeCtx)
	log.Println("Shutdown complete")
}
