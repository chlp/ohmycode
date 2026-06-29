package store

import (
	"testing"
	"time"
)

const (
	testRunnerA = "AAAAAAAAAAAAAAAAAAAAAA"
	testRunnerB = "BBBBBBBBBBBBBBBBBBBBBB"
)

func TestRunnerStore_SetRunner_CreatesNew(t *testing.T) {
	rs := NewRunnerStore()
	r := rs.SetRunner(testRunnerA, false)
	if r == nil {
		t.Fatal("SetRunner returned nil")
	}
	if r.ID != testRunnerA {
		t.Errorf("ID: got %q, want %q", r.ID, testRunnerA)
	}
	if r.IsPublic {
		t.Error("IsPublic should be false")
	}
}

func TestRunnerStore_SetRunner_UpdatesExisting(t *testing.T) {
	rs := NewRunnerStore()
	rs.SetRunner(testRunnerA, false)
	r := rs.SetRunner(testRunnerA, true) // flip to public
	if !r.IsPublic {
		t.Error("SetRunner should update IsPublic on second call")
	}
}

func TestRunnerStore_IsOnline_FreshRunner(t *testing.T) {
	rs := NewRunnerStore()
	rs.SetRunner(testRunnerA, false)
	if !rs.IsOnline(false, testRunnerA) {
		t.Error("freshly-set runner should be online")
	}
}

func TestRunnerStore_IsOnline_UnknownRunner_False(t *testing.T) {
	rs := NewRunnerStore()
	if rs.IsOnline(false, testRunnerA) {
		t.Error("unknown runner should be offline")
	}
}

func TestRunnerStore_IsOnline_PublicRunner(t *testing.T) {
	rs := NewRunnerStore()
	rs.SetRunner(testRunnerA, true)
	if !rs.IsOnline(true, "") {
		t.Error("fresh public runner should be online when isPublic=true")
	}
}

func TestRunnerStore_CountOnline_TwoFresh(t *testing.T) {
	rs := NewRunnerStore()
	rs.SetRunner(testRunnerA, false)
	rs.SetRunner(testRunnerB, false)
	if n := rs.CountOnline(); n != 2 {
		t.Errorf("CountOnline: got %d, want 2", n)
	}
}

func TestRunnerStore_CountOnline_EmptyStore(t *testing.T) {
	rs := NewRunnerStore()
	if n := rs.CountOnline(); n != 0 {
		t.Errorf("CountOnline on empty store: got %d, want 0", n)
	}
}

func TestRunnerStore_GetPublicRunner_ReturnsPublic(t *testing.T) {
	rs := NewRunnerStore()
	rs.SetRunner(testRunnerA, true)
	r := rs.GetPublicRunner()
	if r == nil {
		t.Fatal("GetPublicRunner returned nil for a registered public runner")
	}
	if r.ID != testRunnerA {
		t.Errorf("runner ID: got %q, want %q", r.ID, testRunnerA)
	}
}

func TestRunnerStore_GetPublicRunner_NilWhenNone(t *testing.T) {
	rs := NewRunnerStore()
	if rs.GetPublicRunner() != nil {
		t.Error("GetPublicRunner should return nil when no public runners exist")
	}
}

func TestRunnerStore_GetPublicRunner_SkipsPrivate(t *testing.T) {
	rs := NewRunnerStore()
	rs.SetRunner(testRunnerA, false) // private only
	if rs.GetPublicRunner() != nil {
		t.Error("GetPublicRunner should return nil when only private runners exist")
	}
}

func TestRunnerStore_GetPublicRunner_ReturnsMostRecent(t *testing.T) {
	rs := NewRunnerStore()
	rs.SetRunner(testRunnerA, true)
	time.Sleep(10 * time.Millisecond)
	rs.SetRunner(testRunnerB, true)
	r := rs.GetPublicRunner()
	if r == nil {
		t.Fatal("GetPublicRunner returned nil")
	}
	if r.ID != testRunnerB {
		t.Errorf("expected most recently set runner B, got %q", r.ID)
	}
}

func TestRunnerStore_TouchRunner_AdvancesCheckedAt(t *testing.T) {
	rs := NewRunnerStore()
	rs.SetRunner(testRunnerA, false)
	before := rs.GetRunner(testRunnerA).CheckedAt
	time.Sleep(5 * time.Millisecond)
	rs.TouchRunner(testRunnerA)
	after := rs.GetRunner(testRunnerA).CheckedAt
	if !after.After(before) {
		t.Error("TouchRunner should advance CheckedAt")
	}
}

func TestRunnerStore_TouchRunner_InvalidId_NoOp(t *testing.T) {
	rs := NewRunnerStore()
	rs.TouchRunner("bad-id") // should not panic
}

func TestRunnerStore_GetRunner_ReturnsNilForUnknown(t *testing.T) {
	rs := NewRunnerStore()
	if rs.GetRunner(testRunnerA) != nil {
		t.Error("GetRunner should return nil for unknown runner")
	}
}
