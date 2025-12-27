package store

import (
	"errors"
	"ohmycode_api/internal/model"
	"sync"
)

type FileStore struct {
	filesMu     sync.RWMutex
	files       map[string]*model.File
	fileLocksMu sync.Mutex
	fileLocks   map[string]*sync.Mutex
	db          *Db
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
	var fromDb model.File
	found, err := fs.db.FindOne("files", map[string]interface{}{"_id": fileId}, &fromDb)
	if err != nil {
		return nil, err
	}
	if found && fromDb.ID == fileId {
		fromDb.Persisted = true
		_ = fromDb.Updates()
		fs.filesMu.Lock()
		fs.files[fileId] = &fromDb
		fs.filesMu.Unlock()
		return &fromDb, nil
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

	var fromDb model.File
	found, err := fs.db.FindOne("files", map[string]interface{}{"_id": fileId}, &fromDb)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	if fromDb.ID != fileId {
		return nil, errors.New("problem with fileId")
	}

	fromDb.Persisted = true
	_ = fromDb.Updates()

	fs.filesMu.Lock()
	fs.files[fileId] = &fromDb
	fs.filesMu.Unlock()
	return &fromDb, nil
}

func (fs *FileStore) PersistFile(file *model.File) error {
	fileMu := fs.lockFileMutex(file.ID)
	defer fileMu.Unlock()

	// Persist a stable snapshot to avoid data races with WS/worker.
	snap := file.Snapshot(true)
	doc := model.File{
		ID:                 snap.ID,
		Name:               snap.Name,
		Lang:               snap.Lang,
		Content:            snap.Content,
		ContentUpdatedAt:   snap.ContentUpdatedAt,
		Result:             snap.Result,
		Writer:             snap.Writer,
		UsePublicRunner:    snap.UsePublicRunner,
		RunnerId:           snap.RunnerId,
		Users:              snap.Users,
		UpdatedAt:          snap.UpdatedAt,
		Persisted:          snap.Persisted,
		IsWaitingForResult: snap.IsWaitingForResult,
		IsRunnerOnline:     snap.IsRunnerOnline,
		PersistedAt:        snap.PersistedAt,
	}

	if err := fs.db.ReplaceOneUpsert("files", map[string]interface{}{"_id": doc.ID}, &doc); err != nil {
		return err
	}
	file.SetPersistedAt(snap.UpdatedAt)
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
