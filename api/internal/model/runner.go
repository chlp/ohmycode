package model

import (
	"time"
)

type Runner struct {
	ID        string    `json:"id"`
	IsPublic  bool      `json:"is_public"`
	CheckedAt time.Time `json:"checked_at"`
}
