package api

import (
	"context"
	"log"
	"net/http"
	"ohmycode_api/internal/store"
	"ohmycode_api/pkg/util"
	"strconv"
	"time"
)

type Service struct {
	httpPort         int
	serveClientFiles bool
	useDynamicFiles  bool
	fileStore        *store.FileStore
	runnerStore      *store.RunnerStore
	taskStore        *store.TaskStore
}

func NewService(httpPort int, serveClientFiles, useDynamicFiles bool, fileStore *store.FileStore, runnerStore *store.RunnerStore, taskStore *store.TaskStore) *Service {
	return &Service{
		httpPort:         httpPort,
		serveClientFiles: serveClientFiles,
		useDynamicFiles:  useDynamicFiles,
		fileStore:        fileStore,
		runnerStore:      runnerStore,
		taskStore:        taskStore,
	}
}

func (s *Service) Run() {
	util.Log("API Service started")
	mux := http.NewServeMux()

	mux.HandleFunc("/file", s.handleWsFileConnection)
	mux.HandleFunc("/runner", s.handleWsRunnerConnection)
	if s.serveClientFiles {
		if s.useDynamicFiles {
			serveDynamicFiles(mux)
		} else {
			serveStaticFiles(mux)
		}
	}

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(s.httpPort), corsMiddleware(timerMiddleware(mux))))
}

func timerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), util.RequestStartTimeCtxKey, time.Now())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
