package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

const persistErrDegradedThreshold = 5

type healthResponse struct {
	Status        string `json:"status"`
	Mongo         string `json:"mongo"`
	PersistErrors int    `json:"persist_errors"`
	RunnersOnline int    `json:"runners_online"`
}

func (s *Service) handleHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	mongoStatus := "ok"
	if err := s.fileStore.Ping(ctx); err != nil {
		mongoStatus = "error"
	}

	persistErrors := s.fileStore.PersistErrCount()

	resp := healthResponse{
		Status:        "ok",
		Mongo:         mongoStatus,
		PersistErrors: persistErrors,
		RunnersOnline: s.runnerStore.CountOnline(),
	}

	degraded := mongoStatus != "ok" || persistErrors >= persistErrDegradedThreshold
	w.Header().Set("Content-Type", "application/json")
	if degraded {
		w.WriteHeader(http.StatusServiceUnavailable)
		resp.Status = "degraded"
	}
	_ = json.NewEncoder(w).Encode(resp)
}
