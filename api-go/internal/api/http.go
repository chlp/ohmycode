package api

import (
	"context"
	"log"
	"net/http"
	"ohmycode_api/internal/store"
	"ohmycode_api/pkg/util"
	"time"
)

func Run(store store.Store) {
	s := NewService(store)

	mux := http.NewServeMux()
	mux.HandleFunc("/file/get", s.HandleFileGetUpdateRequest)
	mux.HandleFunc("/file/set_code", s.HandleFileSetCodeRequest)
	log.Fatal(http.ListenAndServe(":8081", requestTimerMiddleware(mux)))
}

func requestTimerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), util.RequestStartTimeCtxKey, time.Now())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
