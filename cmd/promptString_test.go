package cmd

import (
	"strings"
	"testing"
)

func TestBuildPrompt_WithAlias(t *testing.T) {
	result := buildPrompt("myrepo", "feature/auth", "auth", "/home/user/myrepo", " ⌥ ")
	expected := "myrepo ⌥ feature/auth (auth)"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestBuildPrompt_WithoutAlias(t *testing.T) {
	result := buildPrompt("myrepo", "main", "", "/home/user/myrepo", " ⌥ ")
	expected := "myrepo ⌥ main"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestBuildPrompt_AtRepoRoot(t *testing.T) {
	result := buildPrompt("myrepo", "main", "", "/home/user/myrepo", " ⌥ ")
	// No subpath when at repo root
	if strings.Contains(result, "[") {
		t.Errorf("expected no subpath brackets at repo root, got %q", result)
	}
}

func TestBuildPrompt_InSubdirectory(t *testing.T) {
	result := buildPrompt("myrepo", "main", "", "/home/user/myrepo/src/utils", " ⌥ ")
	expected := "myrepo [src/utils] ⌥ main"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestBuildPrompt_NestedSubdirectory(t *testing.T) {
	result := buildPrompt("myrepo", "main", "", "/home/user/myrepo/a/b/c", " ⌥ ")
	expected := "myrepo [a/b/c] ⌥ main"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestBuildPrompt_WindowsSeparator(t *testing.T) {
	result := buildPrompt("myrepo", "main", "m", "/home/user/myrepo", " > ")
	expected := "myrepo > main (m)"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestBuildPrompt_SubdirectoryWithAlias(t *testing.T) {
	result := buildPrompt("myrepo", "feature/x", "fx", "/home/user/myrepo/src", " ⌥ ")
	expected := "myrepo [src] ⌥ feature/x (fx)"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestBuildPrompt_RepoNotInPath(t *testing.T) {
	// Edge case: repo name not found in working directory
	result := buildPrompt("myrepo", "main", "", "/some/other/path", " ⌥ ")
	// Should still produce output without subpath
	expected := "myrepo ⌥ main"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
