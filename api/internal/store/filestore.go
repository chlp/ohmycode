package store

import (
	"errors"
	"ohmycode_api/internal/model"
	"sync"
)

type FileStore struct {
	filesMu    sync.RWMutex
	files      map[string]*model.File
	fileLocksMu sync.Mutex
	fileLocks  map[string]*sync.Mutex
	db         *Db
}

func NewFileStore(dbConfig DBConfig) *FileStore {
	return &FileStore{
		files:     make(map[string]*model.File),
		fileLocks: make(map[string]*sync.Mutex),
		db:        newDb(dbConfig),
	}
}

func (fs *FileStore) GetFileOrCreate(fileId, fileName, lang, content, userId, userName string) (*model.File, error) {
	fileMu := fs.lockFileMutex(fileId)
	defer fileMu.Unlock()

	// fast path: in-memory
	fs.filesMu.RLock()
	if file, ok := fs.files[fileId]; ok {
		fs.filesMu.RUnlock()
		return file, nil
	}
	fs.filesMu.RUnlock()

	// try DB
	filesRaw, err := fs.db.Select("files", map[string]interface{}{"_id": fileId}, &model.File{})
	if err != nil {
		return nil, err
	}
	files, ok := filesRaw.([]model.File)
	if !ok {
		return nil, errors.New("problem with type conversion")
	}
	if len(files) == 1 && files[0].ID == fileId {
		files[0].Persisted = true
		fs.filesMu.Lock()
		fs.files[fileId] = &files[0]
		fs.filesMu.Unlock()
		return &files[0], nil
	}

	// create new
	file := model.NewFile(fileId, fileName, lang, content, userId, userName)
	fs.filesMu.Lock()
	fs.files[fileId] = file
	fs.filesMu.Unlock()
	return file, nil
}

func (fs *FileStore) GetAllFiles() []*model.File {
	fs.filesMu.RLock()
	defer fs.filesMu.RUnlock()
	result := make([]*model.File, 0, len(fs.files))
	for _, f := range fs.files {
		result = append(result, f)
	}
	return result
}

func (fs *FileStore) GetFile(fileId string) (*model.File, error) {
	fileMu := fs.lockFileMutex(fileId)
	defer fileMu.Unlock()

	fs.filesMu.RLock()
	if file, ok := fs.files[fileId]; ok {
		fs.filesMu.RUnlock()
		return file, nil
	}
	fs.filesMu.RUnlock()

	filesRaw, err := fs.db.Select("files", map[string]interface{}{"_id": fileId}, &model.File{})
	if err != nil {
		return nil, err
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
	if files[0].ID != fileId {
		return nil, errors.New("problem with fileId")
	}

	files[0].Persisted = true

	fs.filesMu.Lock()
	fs.files[fileId] = &files[0]
	fs.filesMu.Unlock()
	return &files[0], nil
}

func (fs *FileStore) PersistFile(file *model.File) error {
	fileMu := fs.lockFileMutex(file.ID)
	defer fileMu.Unlock()
	if err := fs.db.Upsert("files", file); err != nil {
		return err
	}
	file.PersistedAt = file.UpdatedAt
	return nil
}

func (fs *FileStore) DeleteFile(fileId string) {
	fileMu := fs.lockFileMutex(fileId)
	defer fileMu.Unlock()
	fs.filesMu.Lock()
	delete(fs.files, fileId)
	fs.filesMu.Unlock()

	// prevent unbounded growth
	fs.fileLocksMu.Lock()
	delete(fs.fileLocks, fileId)
	fs.fileLocksMu.Unlock()
}

func (fs *FileStore) lockFileMutex(fileId string) *sync.Mutex {
	fs.fileLocksMu.Lock()
	defer fs.fileLocksMu.Unlock()
	var fileMutex *sync.Mutex
	var ok bool
	if fileMutex, ok = fs.fileLocks[fileId]; !ok {
		fileMutex = &sync.Mutex{}
		fs.fileLocks[fileId] = fileMutex
	}

	fileMutex.Lock()
	return fileMutex
}
