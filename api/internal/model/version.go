package model

import "time"

type FileVersion struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	FileID    string    `json:"file_id" bson:"file_id"`
	Content   string    `json:"content" bson:"content"`
	Name      string    `json:"name" bson:"name"`
	Lang      string    `json:"lang" bson:"lang"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
