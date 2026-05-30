package store_test

import (
	"sync"
	"testing"

	"ohmycode_api/internal/store"
)

const (
	fsFileA = "AAAAAAAAAAAAAAAAAAAAAA" // valid 22-char base62
	fsFileB = "BBBBBBBBBBBBBBBBBBBBBB"
)

func TestFileStoreInMemory_GetFileOrCreate_CreatesFile(t *testing.T) {
	s := store.NewFileStoreInMemory()

	f, err := s.GetFileOrCreate(fsFileA, "My File", "markdown", "", "u1", "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil file")
	}
	snap := f.Snapshot(false)
	if snap.ID != fsFileA {
		t.Errorf("ID: got %q, want %q", snap.ID, fsFileA)
	}
	if snap.Name != "My File" {
		t.Errorf("Name: got %q, want 'My File'", snap.Name)
	}
}

func TestFileStoreInMemory_GetFileOrCreate_ReturnsSamePointer(t *testing.T) {
	s := store.NewFileStoreInMemory()

	f1, _ := s.GetFileOrCreate(fsFileA, "First", "markdown", "", "u1", "")
	f2, _ := s.GetFileOrCreate(fsFileA, "Second", "markdown", "", "u2", "")

	if f1 != f2 {
		t.Error("same file_id should return the same *File pointer")
	}
	// Name comes from first creation; second call hits the in-memory fast path
	if snap := f1.Snapshot(false); snap.Name != "First" {
		t.Errorf("name: got %q, want 'First'", snap.Name)
	}
}

func TestFileStoreInMemory_GetAllFiles(t *testing.T) {
	s := store.NewFileStoreInMemory()
	s.GetFileOrCreate(fsFileA, "A", "markdown", "", "u1", "")
	s.GetFileOrCreate(fsFileB, "B", "markdown", "", "u2", "")

	if got := len(s.GetAllFiles()); got != 2 {
		t.Errorf("GetAllFiles: got %d, want 2", got)
	}
}

func TestFileStoreInMemory_DeleteFile_AllowsRecreation(t *testing.T) {
	s := store.NewFileStoreInMemory()
	s.GetFileOrCreate(fsFileA, "Original", "markdown", "", "u1", "")
	s.DeleteFile(fsFileA)

	f, _ := s.GetFileOrCreate(fsFileA, "After Delete", "markdown", "", "u2", "")
	if snap := f.Snapshot(false); snap.Name != "After Delete" {
		t.Errorf("after deletion: name=%q, want 'After Delete'", snap.Name)
	}
}

func TestFileStoreInMemory_PersistFile_IsNoopWithoutDB(t *testing.T) {
	s := store.NewFileStoreInMemory()
	f, _ := s.GetFileOrCreate(fsFileA, "Test", "markdown", "", "u1", "")
	if err := s.PersistFile(f); err != nil {
		t.Errorf("PersistFile on in-memory store should be a no-op, got: %v", err)
	}
}

func TestFileStoreInMemory_Concurrent_ReturnsSameFile(t *testing.T) {
	s := store.NewFileStoreInMemory()
	const n = 30

	type result struct {
		id string
	}
	results := make([]result, n)

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			f, _ := s.GetFileOrCreate(fsFileA, "Race Test", "markdown", "", "u1", "")
			if f != nil {
				results[idx] = result{id: f.Snapshot(false).ID}
			}
		}(i)
	}
	wg.Wait()

	for i, r := range results {
		if r.id != fsFileA {
			t.Errorf("goroutine %d: id=%q, want %q", i, r.id, fsFileA)
		}
	}
}
