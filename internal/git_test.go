package internal

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"testing"
)

// --- parseBranchOutput tests (pure function) ---

func TestParseBranchOutput_EmptyString(t *testing.T) {
	result := parseBranchOutput("")
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestParseBranchOutput_SingleBranchNoStar(t *testing.T) {
	result := parseBranchOutput("  main")
	if len(result) != 1 || result[0] != "main" {
		t.Errorf("expected [main], got %v", result)
	}
}

func TestParseBranchOutput_SingleBranchWithStar(t *testing.T) {
	result := parseBranchOutput("* main")
	if len(result) != 1 || result[0] != "main" {
		t.Errorf("expected [main], got %v", result)
	}
}

func TestParseBranchOutput_MultipleBranches(t *testing.T) {
	input := "* main\n  feature-a\n  bugfix-1"
	result := parseBranchOutput(input)
	expected := []string{"main", "feature-a", "bugfix-1"}
	if len(result) != len(expected) {
		t.Fatalf("expected %d branches, got %d: %v", len(expected), len(result), result)
	}
	for i, b := range expected {
		if result[i] != b {
			t.Errorf("index %d: expected %q, got %q", i, b, result[i])
		}
	}
}

func TestParseBranchOutput_WhitespaceOnlyLines(t *testing.T) {
	input := "  \n* main\n  \n"
	result := parseBranchOutput(input)
	if len(result) != 1 || result[0] != "main" {
		t.Errorf("expected [main], got %v", result)
	}
}

func TestParseBranchOutput_BranchWithSlashes(t *testing.T) {
	input := "  feature/auth/oauth"
	result := parseBranchOutput(input)
	if len(result) != 1 || result[0] != "feature/auth/oauth" {
		t.Errorf("expected [feature/auth/oauth], got %v", result)
	}
}

func TestParseBranchOutput_DetachedHead(t *testing.T) {
	input := "* (HEAD detached at abc1234)\n  main\n  feature-x"
	result := parseBranchOutput(input)
	// Detached HEAD should be filtered out
	expected := []string{"main", "feature-x"}
	if len(result) != len(expected) {
		t.Fatalf("expected %d branches (no detached HEAD), got %d: %v", len(expected), len(result), result)
	}
	for i, b := range expected {
		if result[i] != b {
			t.Errorf("index %d: expected %q, got %q", i, b, result[i])
		}
	}
}

func TestParseBranchOutput_OnlyDetachedHead(t *testing.T) {
	input := "* (HEAD detached at abc1234)"
	result := parseBranchOutput(input)
	if len(result) != 0 {
		t.Errorf("expected empty slice for detached-only, got %v", result)
	}
}

func TestParseBranchOutput_MixedSpacing(t *testing.T) {
	input := "  main\n    feature-a\n*   bugfix-1"
	result := parseBranchOutput(input)
	if len(result) != 3 {
		t.Fatalf("expected 3 branches, got %d: %v", len(result), result)
	}
	if result[0] != "main" || result[1] != "feature-a" || result[2] != "bugfix-1" {
		t.Errorf("unexpected result: %v", result)
	}
}

// --- Git integration tests (require real git) ---

// initTestRepo creates a temp git repo with an initial commit and returns
// its path and a cleanup function. Tests should call cleanup via t.Cleanup.
func initTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	cmds := [][]string{
		{"git", "init", dir},
		{"git", "-C", dir, "config", "user.email", "test@test.com"},
		{"git", "-C", dir, "config", "user.name", "Test"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git command %v failed: %s\n%s", args, err, out)
		}
	}

	// Create initial commit (empty tree)
	emptyFile := filepath.Join(dir, ".gitkeep")
	if err := os.WriteFile(emptyFile, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}
	cmds = [][]string{
		{"git", "-C", dir, "add", ".gitkeep"},
		{"git", "-C", dir, "commit", "-m", "initial commit"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git command %v failed: %s\n%s", args, err, out)
		}
	}

	return dir
}

// runGit runs a git command in the given dir.
func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	fullArgs := append([]string{"-C", dir}, args...)
	cmd := exec.Command("git", fullArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %s\n%s", args, err, out)
	}
	return string(out)
}

func TestIsGitRepo_InsideRepo(t *testing.T) {
	dir := initTestRepo(t)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	g := Git{}
	if !g.IsGitRepo() {
		t.Error("expected IsGitRepo() = true inside a git repo")
	}
}

func TestIsGitRepo_OutsideRepo(t *testing.T) {
	dir := t.TempDir() // plain dir, not a git repo
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	g := Git{}
	if g.IsGitRepo() {
		t.Error("expected IsGitRepo() = false outside a git repo")
	}
}

func TestGetCurrentBranch(t *testing.T) {
	dir := initTestRepo(t)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	g := Git{}
	branch, err := g.GetCurrentBranch()
	if err != nil {
		t.Fatalf("GetCurrentBranch failed: %v", err)
	}
	// Default branch could be "main" or "master" depending on git config
	if branch == "" {
		t.Error("expected non-empty branch name")
	}
}

func TestGetCurrentBranch_Caching(t *testing.T) {
	dir := initTestRepo(t)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	g := Git{}
	branch1, _ := g.GetCurrentBranch()
	branch2, _ := g.GetCurrentBranch()
	if branch1 != branch2 {
		t.Errorf("expected cached result, got %q then %q", branch1, branch2)
	}
}

func TestGetRepositoryName(t *testing.T) {
	dir := initTestRepo(t)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	g := Git{}
	name, err := g.GetRepositoryName()
	if err != nil {
		t.Fatalf("GetRepositoryName failed: %v", err)
	}
	expected := filepath.Base(dir)
	if name != expected {
		t.Errorf("expected %q, got %q", expected, name)
	}
}

func TestGetRepositoryName_Caching(t *testing.T) {
	dir := initTestRepo(t)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	g := Git{}
	name1, _ := g.GetRepositoryName()
	name2, _ := g.GetRepositoryName()
	if name1 != name2 {
		t.Errorf("expected cached result, got %q then %q", name1, name2)
	}
}

func TestCreateAndListBranches(t *testing.T) {
	dir := initTestRepo(t)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	g := Git{}

	// Create branches
	if err := g.CreateNewBranch("feature-a"); err != nil {
		t.Fatalf("CreateNewBranch failed: %v", err)
	}
	if err := g.CreateNewBranch("feature-b"); err != nil {
		t.Fatalf("CreateNewBranch failed: %v", err)
	}

	branches, err := g.GetBranches()
	if err != nil {
		t.Fatalf("GetBranches failed: %v", err)
	}

	// Should have at least the default branch + 2 created
	if len(branches) < 3 {
		t.Fatalf("expected at least 3 branches, got %d: %v", len(branches), branches)
	}

	// Verify our branches are present
	branchSet := make(map[string]bool)
	for _, b := range branches {
		branchSet[b] = true
	}
	if !branchSet["feature-a"] {
		t.Error("expected 'feature-a' in branch list")
	}
	if !branchSet["feature-b"] {
		t.Error("expected 'feature-b' in branch list")
	}
}

func TestDeleteBranch(t *testing.T) {
	dir := initTestRepo(t)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	g := Git{}
	if err := g.CreateNewBranch("to-delete"); err != nil {
		t.Fatalf("CreateNewBranch failed: %v", err)
	}

	// Verify it exists
	branches, _ := g.GetBranches()
	found := false
	for _, b := range branches {
		if b == "to-delete" {
			found = true
		}
	}
	if !found {
		t.Fatal("branch 'to-delete' should exist before deletion")
	}

	// Delete it (force=true since it's unmerged)
	if err := g.DeleteBranch("to-delete", true); err != nil {
		t.Fatalf("DeleteBranch failed: %v", err)
	}

	// Verify it's gone
	branches, _ = g.GetBranches()
	for _, b := range branches {
		if b == "to-delete" {
			t.Error("branch 'to-delete' should not exist after deletion")
		}
	}
}

func TestSwitchBranch(t *testing.T) {
	dir := initTestRepo(t)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	g := Git{}
	if err := g.CreateNewBranch("test-branch"); err != nil {
		t.Fatalf("CreateNewBranch failed: %v", err)
	}

	if err := g.SwitchBranch("test-branch"); err != nil {
		t.Fatalf("SwitchBranch failed: %v", err)
	}

	// Need a fresh Git to avoid cached branch
	g2 := Git{}
	current, err := g2.GetCurrentBranch()
	if err != nil {
		t.Fatalf("GetCurrentBranch failed: %v", err)
	}
	if current != "test-branch" {
		t.Errorf("expected current branch 'test-branch', got %q", current)
	}
}

func TestSwitchToNewBranch(t *testing.T) {
	dir := initTestRepo(t)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	g := Git{}
	if err := g.SwitchToNewBranch("new-feature"); err != nil {
		t.Fatalf("SwitchToNewBranch failed: %v", err)
	}

	g2 := Git{}
	current, err := g2.GetCurrentBranch()
	if err != nil {
		t.Fatalf("GetCurrentBranch failed: %v", err)
	}
	if current != "new-feature" {
		t.Errorf("expected current branch 'new-feature', got %q", current)
	}
}

func TestGetBranches_OnlyDefaultBranch(t *testing.T) {
	dir := initTestRepo(t)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	g := Git{}
	branches, err := g.GetBranches()
	if err != nil {
		t.Fatalf("GetBranches failed: %v", err)
	}
	if len(branches) != 1 {
		t.Errorf("expected 1 branch (default), got %d: %v", len(branches), branches)
	}
}

func TestGetBranches_ManyBranches(t *testing.T) {
	dir := initTestRepo(t)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	g := Git{}
	branchNames := []string{"alpha", "beta", "gamma", "delta", "epsilon"}
	for _, name := range branchNames {
		if err := g.CreateNewBranch(name); err != nil {
			t.Fatalf("CreateNewBranch(%s) failed: %v", name, err)
		}
	}

	branches, err := g.GetBranches()
	if err != nil {
		t.Fatalf("GetBranches failed: %v", err)
	}

	// Should have default + 5 created
	if len(branches) != 6 {
		t.Fatalf("expected 6 branches, got %d: %v", len(branches), branches)
	}

	sort.Strings(branches)
	for _, name := range branchNames {
		idx := sort.SearchStrings(branches, name)
		if idx >= len(branches) || branches[idx] != name {
			t.Errorf("expected branch %q in list", name)
		}
	}
}

func TestDeleteBranch_NonExistent(t *testing.T) {
	dir := initTestRepo(t)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	g := Git{}
	err := g.DeleteBranch("does-not-exist", true)
	if err == nil {
		t.Error("expected error when deleting non-existent branch")
	}
}

func TestGetRepositoryPath(t *testing.T) {
	dir := initTestRepo(t)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	g := Git{}
	path, err := g.GetRepositoryPath()
	if err != nil {
		t.Fatalf("GetRepositoryPath failed: %v", err)
	}

	// Resolve symlinks for comparison (macOS /private/var vs /var)
	resolvedDir, _ := filepath.EvalSymlinks(dir)
	resolvedPath, _ := filepath.EvalSymlinks(path)
	if resolvedPath != resolvedDir {
		t.Errorf("expected %q, got %q", resolvedDir, resolvedPath)
	}
}
