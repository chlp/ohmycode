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
	subs  map[chan struct{}]struct{}
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		mutex: &sync.Mutex{},
		tasks: make(map[string]*model.Task),
		subs:  make(map[chan struct{}]struct{}),
	}
}

func (ts *TaskStore) Subscribe() (ch <-chan struct{}, unsubscribe func()) {
	c := make(chan struct{}, 1)
	ts.mutex.Lock()
	ts.subs[c] = struct{}{}
	ts.mutex.Unlock()
	return c, func() {
		ts.mutex.Lock()
		delete(ts.subs, c)
		close(c)
		ts.mutex.Unlock()
	}
}

func (ts *TaskStore) signalLocked() {
	for c := range ts.subs {
		select {
		case c <- struct{}{}:
		default:
		}
	}
}

func (ts *TaskStore) AddTask(file *model.File) {
	ts.mutex.Lock()
	ts.tasks[file.ID] = &model.Task{
		FileId:          file.ID,
		Content:         *file.Content,
		Lang:            file.Lang,
		Hash:            util.OhMySimpleHash(*file.Content),
		RunnerId:        file.RunnerId,
		IsPublic:        file.UsePublicRunner,
		GivenToRunnerAt: time.Time{},
	}
	ts.signalLocked()
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
	ts.signalLocked()
}

const durationToRetryTask = time.Second * 30

func (ts *TaskStore) GetTasksForRunner(runner *model.Runner) []*model.Task {
	tasks := make([]*model.Task, 0)
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	for _, task := range ts.tasks {
		if time.Since(task.GivenToRunnerAt) < durationToRetryTask {
			continue
		}
		if !(runner.IsPublic && task.IsPublic) && runner.ID != task.RunnerId {
			continue
		}
		task.GivenToRunnerAt = time.Now()
		task.RunnerId = runner.ID
		tasks = append(tasks, task)
	}
	return tasks
}
