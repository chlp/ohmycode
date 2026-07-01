package worker

import (
	"context"
	"ohmycode_runner/internal/api"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func newTestApiClient(t *testing.T) *api.Client {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	return api.NewApiClient(ctx, "test-runner", true, "ws://127.0.0.1:1", "")
}

func chdirToTempDir(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(origWd) })
}

func TestMoveTaskWritesRequestFileAtomically(t *testing.T) {
	chdirToTempDir(t)

	td := NewTaskDistributor(newTestApiClient(t), "runner-1", []string{"python3"})

	task := &api.Task{Lang: "python3", Content: "print(1)", Hash: 42}
	if err := td.moveTask(task); err != nil {
		t.Fatalf("moveTask returned error: %v", err)
	}

	requestPath := filepath.Join(getDirForRequests("python3"), strconv.FormatUint(uint64(task.Hash), 10))
	data, err := os.ReadFile(requestPath)
	if err != nil {
		t.Fatalf("expected request file at %s: %v", requestPath, err)
	}
	if string(data) != "print(1)" {
		t.Errorf("request content = %q, want %q", data, "print(1)")
	}
	if _, err := os.Stat(requestPath + ".tmp"); !os.IsNotExist(err) {
		t.Errorf("temp file should have been renamed away, stat err = %v", err)
	}
}

func TestMoveTaskRejectsUnconfiguredLanguage(t *testing.T) {
	chdirToTempDir(t)

	td := NewTaskDistributor(newTestApiClient(t), "runner-1", []string{"python3"})

	task := &api.Task{Lang: "ruby", Content: "puts 1", Hash: 1}
	if err := td.moveTask(task); err == nil {
		t.Fatal("expected error for unconfigured language, got nil")
	}
}
