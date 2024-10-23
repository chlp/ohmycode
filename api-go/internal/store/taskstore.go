package store

import (
	"ohmycode_api/internal/model"
	"ohmycode_api/pkg/util"
	"sync"
	"time"
)

type TaskStore struct {
	mutex *sync.Mutex
	tasks map[string]*model.Task
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		mutex: &sync.Mutex{},
		tasks: make(map[string]*model.Task),
	}
}

func (ts *TaskStore) AddTask(file *model.File) {
	ts.mutex.Lock()
	ts.tasks[file.ID] = &model.Task{
		FileId:                 file.ID,
		Content:                file.Content,
		Lang:                   file.Lang,
		Hash:                   util.OhMySimpleHash(file.Content),
		RunnerId:               file.RunnerId,
		IsPublic:               file.UsePublicRunner,
		GivenToRunnerAt:        time.Time{},
		AcknowledgedByRunnerAt: time.Time{},
	}
	ts.mutex.Unlock()
}

func (ts *TaskStore) GetTask(runnerId, lang string, hash uint32) *model.Task {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	for _, task := range ts.tasks {
		if task.RunnerId == runnerId && task.Lang == lang && task.Hash == hash {
			return task
		}
	}
	return nil
}

func (ts *TaskStore) GetTaskForFile(fileId string) *model.Task {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	for _, task := range ts.tasks {
		if task.FileId == fileId {
			return task
		}
	}
	return nil
}

func (ts *TaskStore) DeleteTask(taskId string) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	delete(ts.tasks, taskId)
}

const durationToRemoveGivenTime = time.Second * 3
const durationToRemoveAcknowledgedTime = time.Second * 30

func (ts *TaskStore) GetTasksForRunner(runnerId string, isPublic bool) []*model.Task {
	tasks := make([]*model.Task, 0)
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	for _, task := range ts.tasks {
		if time.Since(task.GivenToRunnerAt) < durationToRemoveGivenTime || time.Since(task.AcknowledgedByRunnerAt) < durationToRemoveAcknowledgedTime {
			continue
		}
		if !(isPublic && task.IsPublic) && runnerId != task.RunnerId {
			continue
		}
		task.GivenToRunnerAt = time.Now()
		task.AcknowledgedByRunnerAt = time.Time{}
		task.RunnerId = runnerId
		tasks = append(tasks, task)
	}
	return tasks
}
