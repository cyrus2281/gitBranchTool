package internal

import (
	"path/filepath"
	"testing"
)

// --- ParseWorktreeList tests ---

func TestParseWorktreeList_EmptyString(t *testing.T) {
	result := ParseWorktreeList("")
	if len(result) != 0 {
		t.Errorf("expected empty map, got %d entries", len(result))
	}
}

func TestParseWorktreeList_SingleWorktree(t *testing.T) {
	input := "worktree /home/user/repo\nbranch refs/heads/main\n"
	result := ParseWorktreeList(input)

	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result["/home/user/repo"] != "main" {
		t.Errorf("expected branch 'main', got '%s'", result["/home/user/repo"])
	}
}

func TestParseWorktreeList_MultipleWorktrees(t *testing.T) {
	input := "worktree /home/user/repo\nbranch refs/heads/main\n\nworktree /home/user/wt1\nbranch refs/heads/feature-a\n\nworktree /home/user/wt2\nbranch refs/heads/bugfix-1\n"
	result := ParseWorktreeList(input)

	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	expected := map[string]string{
		"/home/user/repo": "main",
		"/home/user/wt1":  "feature-a",
		"/home/user/wt2":  "bugfix-1",
	}
	for path, branch := range expected {
		if result[path] != branch {
			t.Errorf("path %q: expected branch %q, got %q", path, branch, result[path])
		}
	}
}

func TestParseWorktreeList_DetachedHead(t *testing.T) {
	input := "worktree /home/user/repo\nbranch refs/heads/main\n\nworktree /home/user/detached\nHEAD abc123\ndetached\n"
	result := ParseWorktreeList(input)

	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result["/home/user/repo"] != "main" {
		t.Errorf("expected branch 'main', got '%s'", result["/home/user/repo"])
	}
	// Detached worktree should have empty branch name
	if result["/home/user/detached"] != "" {
		t.Errorf("expected empty branch for detached worktree, got '%s'", result["/home/user/detached"])
	}
}

func TestParseWorktreeList_BranchWithSlashes(t *testing.T) {
	input := "worktree /home/user/repo\nbranch refs/heads/feature/auth/oauth\n"
	result := ParseWorktreeList(input)

	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result["/home/user/repo"] != "feature/auth/oauth" {
		t.Errorf("expected 'feature/auth/oauth', got '%s'", result["/home/user/repo"])
	}
}

func TestParseWorktreeList_TrailingNewlines(t *testing.T) {
	input := "worktree /home/user/repo\nbranch refs/heads/main\n\n\n\n"
	result := ParseWorktreeList(input)

	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result["/home/user/repo"] != "main" {
		t.Errorf("expected 'main', got '%s'", result["/home/user/repo"])
	}
}

func TestParseWorktreeList_MixedDetachedAndNormal(t *testing.T) {
	input := "worktree /repo\nbranch refs/heads/main\n\nworktree /wt-detached\nHEAD abc123\ndetached\n\nworktree /wt-feature\nbranch refs/heads/feature-x\n"
	result := ParseWorktreeList(input)

	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	if result["/repo"] != "main" {
		t.Errorf("expected 'main', got '%s'", result["/repo"])
	}
	if result["/wt-detached"] != "" {
		t.Errorf("expected empty for detached, got '%s'", result["/wt-detached"])
	}
	if result["/wt-feature"] != "feature-x" {
		t.Errorf("expected 'feature-x', got '%s'", result["/wt-feature"])
	}
}

func TestParseWorktreeList_BareWorktreeEntries(t *testing.T) {
	// Worktree entries with no branch line (bare repos or special states)
	input := "worktree /bare1\n\nworktree /bare2\n"
	result := ParseWorktreeList(input)

	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result["/bare1"] != "" {
		t.Errorf("expected empty branch for /bare1, got '%s'", result["/bare1"])
	}
	if result["/bare2"] != "" {
		t.Errorf("expected empty branch for /bare2, got '%s'", result["/bare2"])
	}
}

// --- ResolveWorktreePath tests ---

func TestResolveWorktreePath_DefaultTemplate(t *testing.T) {
	result := ResolveWorktreePath("./worktrees/{alias}", "/home/user/repo", "repo", "my-alias", "feature/branch")
	expected := filepath.Join("/home/user/repo", "worktrees", "my-alias")
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestResolveWorktreePath_AllVariables(t *testing.T) {
	result := ResolveWorktreePath("{repository}/{alias}/{branch}", "/home/user/repo", "myrepo", "wt1", "main")
	expected := filepath.Join("/home/user/repo", "myrepo", "wt1", "main")
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestResolveWorktreePath_AbsolutePath(t *testing.T) {
	result := ResolveWorktreePath("/tmp/worktrees/{alias}", "/home/user/repo", "repo", "wt1", "main")
	expected := filepath.Clean("/tmp/worktrees/wt1")
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestResolveWorktreePath_RelativeParentDir(t *testing.T) {
	result := ResolveWorktreePath("../worktrees/{alias}", "/home/user/repo", "repo", "wt1", "main")
	expected := filepath.Clean("/home/user/worktrees/wt1")
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestResolveWorktreePath_NoVariables(t *testing.T) {
	result := ResolveWorktreePath("./fixed-path", "/home/user/repo", "repo", "wt1", "main")
	expected := filepath.Join("/home/user/repo", "fixed-path")
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestResolveWorktreePath_EmptyAlias(t *testing.T) {
	result := ResolveWorktreePath("./worktrees/{alias}", "/home/user/repo", "repo", "", "main")
	expected := filepath.Join("/home/user/repo", "worktrees")
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestResolveWorktreePath_BranchWithSlashes(t *testing.T) {
	result := ResolveWorktreePath("/tmp/{branch}", "/repo", "repo", "wt1", "feature/auth")
	expected := filepath.Clean("/tmp/feature/auth")
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

// --- GetWorktreePathForBranch tests ---

func TestGetWorktreePathForBranch_Found(t *testing.T) {
	m := map[string]string{
		"/path/to/main": "main",
	}
	result := GetWorktreePathForBranch(m, "main")
	if result != "/path/to/main" {
		t.Errorf("expected '/path/to/main', got '%s'", result)
	}
}

func TestGetWorktreePathForBranch_FoundInMultiple(t *testing.T) {
	m := map[string]string{
		"/path/to/main":    "main",
		"/path/to/feature": "feature-x",
		"/path/to/bugfix":  "bugfix-1",
	}
	result := GetWorktreePathForBranch(m, "feature-x")
	if result != "/path/to/feature" {
		t.Errorf("expected '/path/to/feature', got '%s'", result)
	}
}

func TestGetWorktreePathForBranch_NotFound(t *testing.T) {
	m := map[string]string{
		"/path/to/main": "main",
	}
	result := GetWorktreePathForBranch(m, "nonexistent")
	if result != "" {
		t.Errorf("expected empty string, got '%s'", result)
	}
}

func TestGetWorktreePathForBranch_EmptyMap(t *testing.T) {
	m := map[string]string{}
	result := GetWorktreePathForBranch(m, "main")
	if result != "" {
		t.Errorf("expected empty string, got '%s'", result)
	}
}

// --- Worktree.String() tests ---

func TestWorktreeString(t *testing.T) {
	wt := Worktree{Alias: "wt1", Path: "/path/to/wt", Note: "my note"}
	s := wt.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
	// Verify all fields are present in the output
	if !containsStr(s, "/path/to/wt") || !containsStr(s, "wt1") || !containsStr(s, "my note") {
		t.Errorf("expected path, alias, and note in output, got: %s", s)
	}
}

func TestWorktreeStringWithBranch(t *testing.T) {
	wt := Worktree{Alias: "wt1", Path: "/path/to/wt", Note: "my note"}
	s := wt.StringWithBranch("feature-x", "fx")
	if s == "" {
		t.Error("expected non-empty string")
	}
	if !containsStr(s, "feature-x") || !containsStr(s, "fx") {
		t.Errorf("expected branch and alias in output, got: %s", s)
	}
}

// helper
func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && findSubstr(s, substr))
}

func findSubstr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
