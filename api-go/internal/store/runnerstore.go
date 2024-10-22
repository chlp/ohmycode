package store

import (
	"ohmycode_api/internal/model"
	"ohmycode_api/pkg/util"
	"sync"
	"time"
)

type RunnerStore struct {
	mutex   *sync.RWMutex
	runners map[string]*model.Runner
}

func NewRunnerStore() *RunnerStore {
	return &RunnerStore{
		mutex:   &sync.RWMutex{},
		runners: make(map[string]*model.Runner),
	}
}

func (s *RunnerStore) SetRunner(id string, isPublic bool) {
	if !util.IsUuid(id) {
		return
	}

	s.mutex.RLock()
	if runner, ok := s.runners[id]; ok {
		runner.IsPublic = isPublic
		runner.CheckedAt = time.Now()
		s.mutex.RUnlock()
		return
	}
	s.mutex.RUnlock()

	s.mutex.Lock()
	s.runners[id] = &model.Runner{
		ID:        id,
		IsPublic:  isPublic,
		CheckedAt: time.Now(),
	}
	s.mutex.Unlock()
}

const durationIsActiveFromLastUpdate = 5 * time.Second

func (s *RunnerStore) IsOnline(isPublic bool, runnerId string) bool {
	if isPublic {
		runner := s.GetPublicRunner()
		if runner != nil {
			return time.Since(runner.CheckedAt) < durationIsActiveFromLastUpdate
		}
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if runner, ok := s.runners[runnerId]; ok {
		return time.Since(runner.CheckedAt) < durationIsActiveFromLastUpdate
	}

	return false
}

func (s *RunnerStore) GetPublicRunner() *model.Runner {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	var result *model.Runner
	for _, runner := range s.runners {
		if runner.IsPublic {
			if result == nil || result.CheckedAt.Before(runner.CheckedAt) {
				result = runner
			}
		}
	}
	return result
}
