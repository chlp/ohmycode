package api

import (
	"errors"
	"ohmycode_api/internal/model"
	"ohmycode_api/internal/store"
	"sync"
	"time"
)

type Service struct {
	files   map[string]*model.File
	runners map[string]*model.Runner
	db      store.Db
	mutex   sync.Mutex
}

func NewService(db store.Db) Service {
	return Service{
		files:   make(map[string]*model.File),
		runners: make(map[string]*model.Runner),
		db:      db,
		mutex:   sync.Mutex{},
	}
}

func (s *Service) GetFile(id string) (*model.File, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if file, ok := s.files[id]; ok {
		return file, nil
	}
	var files []model.File
	if err := s.db.Select("files", map[string]interface{}{"_id": id}, files); err != nil {
		return nil, err
	}
	if len(files) == 0 || files[0].ID != id {
		return nil, errors.New("id not found")
	}
	s.files[id] = &files[0]
	return &files[0], nil
}

func (s *Service) GetNewestPublicRunnerCheckedAt() *time.Time {
	var t *time.Time
	for _, runner := range s.runners {
		if runner.IsPublic {
			if t == nil || t.Before(runner.CheckedAt) {
				t = &runner.CheckedAt
			}
		}
	}
	return t
}
