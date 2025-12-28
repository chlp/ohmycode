package model

import (
	"strings"
	"testing"
)

func TestFile_StartRunWithContent_SetsWaitingAndUpdatesContent(t *testing.T) {
	f := NewFile("00000000-0000-0000-0000-000000000000", "t", "php", "<?php echo 1;", "11111111-1111-1111-1111-111111111111", "u")

	// Pretend we already have an old result from a previous run.
	if err := f.SetResult("old result"); err != nil {
		t.Fatalf("SetResult: %v", err)
	}

	const appID = "22222222-2222-2222-2222-222222222222"
	const newContent = "<?php\necho 2;\n"

	if err := f.StartRunWithContent(newContent, appID); err != nil {
		t.Fatalf("StartRunWithContent: %v", err)
	}

	s := f.Snapshot(true)
	if s.Content == nil || *s.Content != newContent {
		t.Fatalf("content mismatch: got=%v", s.Content)
	}
	if s.Writer != appID {
		t.Fatalf("writer mismatch: got=%q want=%q", s.Writer, appID)
	}
	if !s.IsWaitingForResult {
		t.Fatalf("expected IsWaitingForResult=true")
	}
	if strings.Contains(s.Result, "old result") {
		t.Fatalf("result should not keep old value while waiting, got=%q", s.Result)
	}
	if !strings.HasPrefix(s.Result, "Started execution at ") {
		t.Fatalf("unexpected result prefix, got=%q", s.Result)
	}
}


