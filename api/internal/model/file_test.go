package model

import (
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	testFileID = "test-file-id"
	testUserID = "test-user-id"
	testAppID  = "test-app-id"
	testAppID2 = "test-app-id-2"
)

// --- SetContent ---

func TestFile_SetContent_UpdatesFields(t *testing.T) {
	f := NewFile(testFileID, "name", "markdown", "original", testUserID, "u")
	before := f.Snapshot(true)

	if err := f.SetContent("updated", testAppID); err != nil {
		t.Fatalf("SetContent: %v", err)
	}

	s := f.Snapshot(true)
	if s.Content == nil || *s.Content != "updated" {
		t.Errorf("content: got %v, want 'updated'", s.Content)
	}
	if s.Writer != testAppID {
		t.Errorf("writer: got %q, want %q", s.Writer, testAppID)
	}
	if !s.UpdatedAt.After(before.UpdatedAt) {
		t.Error("UpdatedAt should advance after SetContent")
	}
}

func TestFile_SetContent_RejectsOtherWriter(t *testing.T) {
	f := NewFile(testFileID, "name", "markdown", "", testUserID, "u")
	// First writer claims the file
	if err := f.SetContent("v1", testAppID); err != nil {
		t.Fatalf("first SetContent: %v", err)
	}
	// Different app should be rejected
	if err := f.SetContent("v2", testAppID2); err == nil {
		t.Error("expected error when another app tries to write a locked file")
	}
}

func TestFile_SetContent_SameWriterSucceeds(t *testing.T) {
	f := NewFile(testFileID, "name", "markdown", "", testUserID, "u")
	_ = f.SetContent("v1", testAppID)
	if err := f.SetContent("v2", testAppID); err != nil {
		t.Errorf("same writer should succeed on second write: %v", err)
	}
}

func TestFile_SetContent_RejectsTooLong(t *testing.T) {
	f := NewFile(testFileID, "name", "markdown", "", testUserID, "u")
	big := strings.Repeat("x", contentMaxLength+1)
	if err := f.SetContent(big, testAppID); err == nil {
		t.Error("expected error for content exceeding max length")
	}
}

func TestFile_SetContent_NotifiesSubscribers(t *testing.T) {
	f := NewFile(testFileID, "name", "markdown", "", testUserID, "u")
	ch := f.SubscribeUpdates()
	defer f.UnsubscribeUpdates(ch)

	_ = f.SetContent("hello", testAppID)

	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Error("timeout: subscriber did not receive update signal from SetContent")
	}
}

// --- SetLang ---

func TestFile_SetLang_ValidLang(t *testing.T) {
	f := NewFile(testFileID, "name", "markdown", "", testUserID, "u")
	if !f.SetLang("go") {
		t.Error("expected SetLang to return true for valid lang 'go'")
	}
	if s := f.Snapshot(false); s.Lang != "go" {
		t.Errorf("lang not updated: got %q", s.Lang)
	}
}

func TestFile_SetLang_InvalidLangRejected(t *testing.T) {
	f := NewFile(testFileID, "name", "markdown", "", testUserID, "u")
	if f.SetLang("cobol") {
		t.Error("expected SetLang to return false for unknown lang 'cobol'")
	}
	if s := f.Snapshot(false); s.Lang != "markdown" {
		t.Errorf("lang should not change on invalid input, got %q", s.Lang)
	}
}

// --- Snapshot ---

func TestFile_Snapshot_IsDeepCopy(t *testing.T) {
	f := NewFile(testFileID, "name", "markdown", "original", testUserID, "u")
	s := f.Snapshot(true)

	// Modify file — snapshot should be unaffected
	_ = f.SetContent("modified", testAppID)

	if s.Content == nil || *s.Content != "original" {
		t.Error("snapshot content should be independent from the original file")
	}

	origUserCount := len(s.Users)
	f.TouchByUser("another-user-id-12345", "Bob")
	if len(s.Users) != origUserCount {
		t.Error("snapshot users slice should be independent from the original file")
	}
}

// --- Subscribe / Unsubscribe ---

func TestFile_UnsubscribeUpdates_StopsSignals(t *testing.T) {
	f := NewFile(testFileID, "name", "markdown", "", testUserID, "u")
	ch := f.SubscribeUpdates()
	f.UnsubscribeUpdates(ch)

	// Channel is closed after unsubscribe; writing to file should not panic
	_ = f.SetContent("after-unsub", testAppID)

	// Closed channel is immediately readable (returns zero value)
	select {
	case _, ok := <-ch:
		if ok {
			t.Error("channel should be closed after unsubscribe")
		}
	default:
		t.Error("closed channel should be immediately readable")
	}
}

// --- Concurrency ---

func TestFile_ConcurrentSetContent_IsRaceFree(t *testing.T) {
	f := NewFile(testFileID, "name", "markdown", "", testUserID, "u")

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = f.SetContent("content", testAppID) // same appID — all succeed
		}()
	}
	wg.Wait()

	s := f.Snapshot(true)
	if s.Content == nil || *s.Content != "content" {
		t.Error("unexpected content after concurrent SetContent calls")
	}
}

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


