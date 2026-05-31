package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type healthResponse struct {
	Status        string `json:"status"`
	Mongo         string `json:"mongo"`
	RunnersOnline int    `json:"runners_online"`
}

func (s *Service) handleHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	mongoStatus := "ok"
	if err := s.fileStore.Ping(ctx); err != nil {
		mongoStatus = "error"
	}

	resp := healthResponse{
		Status:        "ok",
		Mongo:         mongoStatus,
		RunnersOnline: s.runnerStore.CountOnline(),
	}

	w.Header().Set("Content-Type", "application/json")
	if mongoStatus != "ok" {
		w.WriteHeader(http.StatusServiceUnavailable)
		resp.Status = "degraded"
	}
	_ = json.NewEncoder(w).Encode(resp)
}
