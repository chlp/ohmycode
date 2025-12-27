package api

import (
	"ohmycode_api/internal/model"
	"time"
)

// fileDTO is a JSON snapshot of model.File without sync primitives.
// It intentionally keeps the JSON shape expected by the frontend.
type fileDTO struct {
	ID               string      `json:"id"`
	Name             string      `json:"name"`
	Lang             string      `json:"lang"`
	Content          *string     `json:"content,omitempty"`
	ContentUpdatedAt time.Time   `json:"content_updated_at"`
	Result           string      `json:"result"`
	WriterID         string      `json:"writer_id"`
	Runner           string      `json:"runner"`
	Users            []model.User `json:"users"`
	UpdatedAt        time.Time   `json:"updated_at"`
	Persisted        bool        `json:"persisted"`
	IsWaitingForResult bool      `json:"is_waiting_for_result"`
	IsRunnerOnline     bool      `json:"is_runner_online"`
}

func toFileDTO(f *model.File, includeContent bool) fileDTO {
	dto := fileDTO{
		ID:                f.ID,
		Name:              f.Name,
		Lang:              f.Lang,
		ContentUpdatedAt:  f.ContentUpdatedAt,
		Result:            f.Result,
		WriterID:          f.Writer,
		Runner:            f.RunnerId,
		Users:             f.Users,
		UpdatedAt:         f.UpdatedAt,
		Persisted:         f.Persisted,
		IsWaitingForResult: f.IsWaitingForResult,
		IsRunnerOnline:     f.IsRunnerOnline,
	}
	if includeContent {
		dto.Content = f.Content
	}
	return dto
}


