package model

import "time"

type Task struct {
	FileId   string
	Content  string `json:"content"`
	Lang     string `json:"lang"`
	Hash     uint32 `json:"hash"`
	RunnerId string
	IsPublic bool

	GivenToRunnerAt        time.Time
	AcknowledgedByRunnerAt time.Time
}
