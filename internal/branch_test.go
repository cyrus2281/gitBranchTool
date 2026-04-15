package internal

import (
	"strings"
	"testing"
)

func TestBranchString_AllFields(t *testing.T) {
	b := Branch{Name: "feature/auth", Alias: "auth", Note: "Working on auth"}
	s := b.String()
	if !strings.Contains(s, "feature/auth") {
		t.Errorf("expected name in output, got: %s", s)
	}
	if !strings.Contains(s, "auth") {
		t.Errorf("expected alias in output, got: %s", s)
	}
	if !strings.Contains(s, "Working on auth") {
		t.Errorf("expected note in output, got: %s", s)
	}
}

func TestBranchString_EmptyFields(t *testing.T) {
	b := Branch{Name: "", Alias: "", Note: ""}
	s := b.String()
	// Should not panic and should produce formatted output
	if s == "" {
		t.Error("expected non-empty formatted string even with empty fields")
	}
}

func TestBranchString_LongValues(t *testing.T) {
	longName := strings.Repeat("a", 50)
	longAlias := strings.Repeat("b", 50)
	longNote := strings.Repeat("c", 50)
	b := Branch{Name: longName, Alias: longAlias, Note: longNote}
	s := b.String()
	if !strings.Contains(s, longName) {
		t.Errorf("expected long name in output")
	}
	if !strings.Contains(s, longAlias) {
		t.Errorf("expected long alias in output")
	}
	if !strings.Contains(s, longNote) {
		t.Errorf("expected long note in output")
	}
}
