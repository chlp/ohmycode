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

func (rs *RunnerStore) GetRunner(id string) *model.Runner {
	if !util.IsUuid(id) {
		return nil
	}
	rs.mutex.RLock()
	if runner, ok := rs.runners[id]; ok {
		rs.mutex.RUnlock()
		return runner
	}
	rs.mutex.RUnlock()
	return nil
}

func (rs *RunnerStore) SetRunner(id string, isPublic bool) *model.Runner {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	if runner, ok := rs.runners[id]; ok {
		runner.IsPublic = isPublic
		runner.CheckedAt = time.Now()
		return runner
	}
	runner := &model.Runner{
		ID:        id,
		IsPublic:  isPublic,
		CheckedAt: time.Now(),
	}
	rs.runners[id] = runner
	return runner
}

const durationIsActiveFromLastUpdate = 5 * time.Second

func (rs *RunnerStore) TouchRunner(id string) {
	if !util.IsUuid(id) {
		return
	}
	rs.mutex.Lock()
	if runner, ok := rs.runners[id]; ok {
		runner.CheckedAt = time.Now()
	}
	rs.mutex.Unlock()
}

func (rs *RunnerStore) IsOnline(isPublic bool, runnerId string) bool {
	if isPublic {
		runner := rs.GetPublicRunner()
		if runner != nil {
			return time.Since(runner.CheckedAt) < durationIsActiveFromLastUpdate
		}
	}

	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	if runner, ok := rs.runners[runnerId]; ok {
		return time.Since(runner.CheckedAt) < durationIsActiveFromLastUpdate
	}

	return false
}

func (rs *RunnerStore) GetPublicRunner() *model.Runner {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	var result *model.Runner
	for _, runner := range rs.runners {
		if runner.IsPublic {
			if result == nil || result.CheckedAt.Before(runner.CheckedAt) {
				result = runner
			}
		}
	}
	return result
}
