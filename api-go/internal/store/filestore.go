package store

import (
	"errors"
	"ohmycode_api/internal/model"
	"sync"
)

type Store struct {
	mutex      *sync.Mutex
	files      map[string]*model.File
	filesMutex map[string]*sync.Mutex
	runners    map[string]*model.Runner
	db         *Db
}

func NewStore(dbConfig DBConfig) Store {
	return Store{
		mutex:      &sync.Mutex{},
		files:      make(map[string]*model.File),
		filesMutex: make(map[string]*sync.Mutex),
		runners:    make(map[string]*model.Runner),
		db:         newDb(dbConfig),
	}
}

func (s *Store) GetFile(id string) (*model.File, error) {
	s.mutex.Lock()
	var fileMutex *sync.Mutex
	var ok bool
	if fileMutex, ok = s.filesMutex[id]; !ok {
		fileMutex = &sync.Mutex{}
		s.filesMutex[id] = fileMutex
	}
	fileMutex.Lock()
	defer fileMutex.Unlock()
	s.mutex.Unlock()

	if file, ok := s.files[id]; ok {
		return file, nil
	}

	filesRaw, err := s.db.Select("files", map[string]interface{}{"_id": id}, &model.File{})
	if err != nil {
		panic(err)
	}
	files, ok := filesRaw.([]model.File)
	if !ok {
		return nil, errors.New("problem with type conversion")
	}
	if len(files) == 0 {
		return nil, nil
	}
	if len(files) != 1 {
		return nil, errors.New("wrong count")
	}
	if files[0].ID != id {
		return nil, errors.New("problem with id")
	}

	s.mutex.Lock()
	s.files[id] = &files[0]
	s.mutex.Unlock()
	return &files[0], nil
}
