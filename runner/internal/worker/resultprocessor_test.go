package worker

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResultProcessorKeepsResultFileWhenSetResultFails(t *testing.T) {
	chdirToTempDir(t)

	rp := NewResultProcessor(newTestApiClient(t), "runner-1", "python3")
	resultsDir := getDirForResults("python3")
	resultPath := filepath.Join(resultsDir, "123")
	if err := os.WriteFile(resultPath, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	processed, err := rp.Process()
	if processed != 1 {
		t.Fatalf("processed = %d, want 1", processed)
	}
	// No live API connection is available in this test, so SetResult must fail
	// and the result file must be left in place for the next retry.
	if err == nil {
		t.Fatal("expected an error since there is no live API connection")
	}
	if _, statErr := os.Stat(resultPath); statErr != nil {
		t.Errorf("expected result file to remain after failed SetResult, stat err = %v", statErr)
	}
}

func TestResultProcessorRemovesInvalidFilename(t *testing.T) {
	chdirToTempDir(t)

	rp := NewResultProcessor(newTestApiClient(t), "runner-1", "python3")
	resultsDir := getDirForResults("python3")
	badPath := filepath.Join(resultsDir, "not-a-number")
	if err := os.WriteFile(badPath, []byte("junk"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := rp.Process(); err == nil {
		t.Fatal("expected error for invalid filename")
	}
	if _, err := os.Stat(badPath); !os.IsNotExist(err) {
		t.Errorf("invalid result file should have been removed, stat err = %v", err)
	}
}

func TestResultProcessorIgnoresDotfiles(t *testing.T) {
	chdirToTempDir(t)

	rp := NewResultProcessor(newTestApiClient(t), "runner-1", "python3")
	resultsDir := getDirForResults("python3")
	if err := os.WriteFile(filepath.Join(resultsDir, ".DS_Store"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	processed, err := rp.Process()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if processed != 0 {
		t.Fatalf("processed = %d, want 0 (dotfile should be skipped)", processed)
	}
}
