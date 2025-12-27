package api

import (
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
	wsAllowedOrigins []string
	fileStore        *store.FileStore
	runnerStore      *store.RunnerStore
	taskStore        *store.TaskStore
}

func NewService(httpPort int, serveClientFiles, useDynamicFiles bool, wsAllowedOrigins []string, fileStore *store.FileStore, runnerStore *store.RunnerStore, taskStore *store.TaskStore) *Service {
	return &Service{
		httpPort:         httpPort,
		serveClientFiles: serveClientFiles,
		useDynamicFiles:  useDynamicFiles,
		wsAllowedOrigins: wsAllowedOrigins,
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

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(s.httpPort), headersMiddleware(timerMiddleware(mux))))
}

func timerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := util.WithRequestStartTime(r.Context(), time.Now())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func headersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Cache-Control", "public, max-age=86400, stale-while-revalidate=604800, stale-if-error=604800")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
