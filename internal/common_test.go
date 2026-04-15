package internal

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// --- contains tests (from completion.go) ---

func TestContains_Found(t *testing.T) {
	slice := []string{"alpha", "beta", "gamma"}
	if !contains(slice, "alpha") {
		t.Error("expected true for first element")
	}
	if !contains(slice, "beta") {
		t.Error("expected true for middle element")
	}
	if !contains(slice, "gamma") {
		t.Error("expected true for last element")
	}
}

func TestContains_NotFound(t *testing.T) {
	slice := []string{"alpha", "beta", "gamma"}
	if contains(slice, "delta") {
		t.Error("expected false for missing element")
	}
}

func TestContains_EmptySlice(t *testing.T) {
	if contains([]string{}, "anything") {
		t.Error("expected false for empty slice")
	}
}

func TestContains_EmptyString(t *testing.T) {
	slice := []string{"alpha", "", "gamma"}
	if !contains(slice, "") {
		t.Error("expected true for empty string in slice")
	}
}

// --- PrintTable tests ---

// captureStdout captures stdout output during fn execution.
func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestPrintTable_SingleColumnSingleRow(t *testing.T) {
	output := captureStdout(func() {
		PrintTable([]string{"Name"}, [][]string{{"alpha"}})
	})

	if !strings.Contains(output, "Name") {
		t.Error("expected header 'Name' in output")
	}
	if !strings.Contains(output, "---") {
		t.Error("expected separator line in output")
	}
	if !strings.Contains(output, "0)") {
		t.Error("expected row index '0)' in output")
	}
	if !strings.Contains(output, "alpha") {
		t.Error("expected 'alpha' in output")
	}
}

func TestPrintTable_MultipleColumnsAndRows(t *testing.T) {
	headers := []string{"Name", "Alias", "Note"}
	rows := [][]string{
		{"feature/auth", "auth", "working on auth"},
		{"bugfix/123", "fix", "critical bug"},
	}
	output := captureStdout(func() {
		PrintTable(headers, rows)
	})

	for _, h := range headers {
		if !strings.Contains(output, h) {
			t.Errorf("expected header %q in output", h)
		}
	}
	if !strings.Contains(output, "0)") || !strings.Contains(output, "1)") {
		t.Error("expected row indices 0) and 1)")
	}
	if !strings.Contains(output, "feature/auth") || !strings.Contains(output, "bugfix/123") {
		t.Error("expected row data in output")
	}
}

func TestPrintTable_WideValues(t *testing.T) {
	headers := []string{"Name"}
	longValue := strings.Repeat("x", 60)
	rows := [][]string{{longValue}}

	output := captureStdout(func() {
		PrintTable(headers, rows)
	})

	if !strings.Contains(output, longValue) {
		t.Error("expected long value in output")
	}
	// The separator should be at least as wide as the value
	lines := strings.Split(output, "\n")
	separatorFound := false
	for _, line := range lines {
		if strings.HasPrefix(line, "---") {
			separatorFound = true
			if len(line) < 60 {
				t.Errorf("separator too short for wide content: %d chars", len(line))
			}
		}
	}
	if !separatorFound {
		t.Error("expected separator line")
	}
}

func TestPrintTable_TenPlusRows(t *testing.T) {
	headers := []string{"Item"}
	rows := make([][]string, 12)
	for i := range rows {
		rows[i] = []string{"item"}
	}

	output := captureStdout(func() {
		PrintTable(headers, rows)
	})

	// Row 10 and 11 should have wider index (e.g., "10) " vs "0) ")
	if !strings.Contains(output, "10)") || !strings.Contains(output, "11)") {
		t.Error("expected indices 10) and 11) in output")
	}
}

func TestPrintTable_NoRows(t *testing.T) {
	output := captureStdout(func() {
		PrintTable([]string{"Name", "Value"}, [][]string{})
	})

	if !strings.Contains(output, "Name") {
		t.Error("expected header even with no rows")
	}
	if !strings.Contains(output, "---") {
		t.Error("expected separator even with no rows")
	}
	// Should not contain any row indices
	if strings.Contains(output, "0)") {
		t.Error("should not have row indices with no data")
	}
}
