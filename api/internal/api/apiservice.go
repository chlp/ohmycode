package api

import (
	"context"
	"errors"
	"net/http"
	"ohmycode_api/internal/store"
	"ohmycode_api/pkg/util"
	"strconv"
	"strings"
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

func (s *Service) Run(ctx context.Context) error {
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

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(s.httpPort),
		Handler: headersMiddleware(timerMiddleware(mux)),
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
		return nil
	case err := <-errCh:
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

func timerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := util.WithRequestStartTime(r.Context(), time.Now())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func isWebSocketUpgrade(r *http.Request) bool {
	if !strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
		return false
	}
	for _, v := range r.Header.Values("Connection") {
		if strings.Contains(strings.ToLower(v), "upgrade") {
			return true
		}
	}
	return false
}

func headersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		// Don't cache WS endpoints / upgrade requests.
		if r.URL.Path == "/file" || r.URL.Path == "/runner" || isWebSocketUpgrade(r) {
			w.Header().Set("Cache-Control", "no-store")
		} else {
			w.Header().Set("Cache-Control", "public, max-age=86400, stale-while-revalidate=604800, stale-if-error=604800")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
