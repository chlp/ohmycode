package store

import "testing"

func TestDiffPreviewNoChange(t *testing.T) {
	if got := diffPreview("same\ncontent", "same\ncontent"); got != "" {
		t.Errorf("diffPreview(same) = %q, want empty", got)
	}
}

func TestDiffPreviewAddedLine(t *testing.T) {
	got := diffPreview("line one", "line one\nline two")
	if got != "+ line two" {
		t.Errorf("diffPreview(added) = %q, want %q", got, "+ line two")
	}
}

func TestDiffPreviewRemovedLine(t *testing.T) {
	got := diffPreview("line one\nline two", "line one")
	if got != "- line two" {
		t.Errorf("diffPreview(removed) = %q, want %q", got, "- line two")
	}
}

func TestDiffPreviewPositionallyChangedLine(t *testing.T) {
	// Same set of lines, different order at position 0 — no add/remove, falls back to "~".
	got := diffPreview("alpha\nbeta", "beta\nalpha")
	if got != "~ beta" {
		t.Errorf("diffPreview(reordered) = %q, want %q", got, "~ beta")
	}
}

func TestDiffPreviewEmptyOldContent(t *testing.T) {
	got := diffPreview("", "first line")
	if got != "+ first line" {
		t.Errorf("diffPreview(empty old) = %q, want %q", got, "+ first line")
	}
}

func TestTruncLineLeavesShortLinesAlone(t *testing.T) {
	short := "+ short line"
	if got := truncLine(short); got != short {
		t.Errorf("truncLine(short) = %q, want unchanged %q", got, short)
	}
}

func TestTruncLineTruncatesLongLinesTo62Runes(t *testing.T) {
	long := "+ " + stringOfLength(80, 'a')
	got := truncLine(long)
	gotRunes := []rune(got)
	if len(gotRunes) != 63 { // 62 chars + ellipsis
		t.Fatalf("truncLine(long) length = %d, want 63", len(gotRunes))
	}
	if gotRunes[62] != '…' {
		t.Errorf("truncLine(long) last rune = %q, want ellipsis", gotRunes[62])
	}
}

func TestSplitTrimmedLinesDropsBlankAndTrimsWhitespace(t *testing.T) {
	lines := splitTrimmedLines("  first  \n\n\tsecond\t\n   \nthird")
	want := []string{"first", "second", "third"}
	if len(lines) != len(want) {
		t.Fatalf("splitTrimmedLines returned %d lines, want %d: %v", len(lines), len(want), lines)
	}
	for i, l := range lines {
		if l != want[i] {
			t.Errorf("line %d = %q, want %q", i, l, want[i])
		}
	}
}

func stringOfLength(n int, r rune) string {
	runes := make([]rune, n)
	for i := range runes {
		runes[i] = r
	}
	return string(runes)
}
