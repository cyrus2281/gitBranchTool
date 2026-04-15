package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cyrus2281/gitBranchTool/internal"
)

func TestExecuteWorktreeCreate_NoBranch(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	executeWorktreeCreate(&git, store, "my-wt", "", "note ",
		"./worktrees/{alias}", repoDir, filepath.Base(repoDir))

	// Should create branch with alias name and auto-register both
	if !store.BranchExists(internal.Branch{Name: "my-wt"}) {
		t.Error("branch should be auto-registered with alias as name")
	}
	if !store.WorktreeAliasExists("my-wt") {
		t.Error("worktree should be registered")
	}

	// Worktree dir should exist
	wtPath := filepath.Join(repoDir, "worktrees", "my-wt")
	if _, err := os.Stat(wtPath); os.IsNotExist(err) {
		t.Error("worktree directory should exist")
	}
}

func TestExecuteWorktreeCreate_WithRegisteredBranch(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	// Create and register a branch
	git.CreateNewBranch("feature/existing")
	store.AddBranch(internal.Branch{Name: "feature/existing", Alias: "ex", Note: ""})

	executeWorktreeCreate(&git, store, "wt-ex", "ex", "",
		"./worktrees/{alias}", repoDir, filepath.Base(repoDir))

	// Should resolve branch name from alias
	if !store.WorktreeAliasExists("wt-ex") {
		t.Error("worktree should be registered")
	}

	// Branch should already be registered (not duplicated)
	branches := store.GetBranches()
	count := 0
	for _, b := range branches {
		if b.Name == "feature/existing" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1 registration of branch, got %d", count)
	}
}

func TestExecuteWorktreeCreate_WithUnregisteredBranch(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	git.CreateNewBranch("feature/unreg-branch")

	executeWorktreeCreate(&git, store, "wt-unreg", "feature/unreg-branch", "",
		"./worktrees/{alias}", repoDir, filepath.Base(repoDir))

	// Worktree should be registered
	if !store.WorktreeAliasExists("wt-unreg") {
		t.Error("worktree should be registered")
	}

	// Branch should NOT be auto-registered (unregistered branch arg)
	if len(store.GetBranches()) != 0 {
		t.Error("branch should NOT be auto-registered for unregistered branch arg")
	}
}

func TestExecuteWorktreeCreate_NewBranchCreated(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	// Pass a branch name that doesn't exist — should create it
	executeWorktreeCreate(&git, store, "wt-new", "brand-new-branch", "",
		"./worktrees/{alias}", repoDir, filepath.Base(repoDir))

	// Worktree should exist
	if !store.WorktreeAliasExists("wt-new") {
		t.Error("worktree should be registered")
	}

	// Branch should NOT be auto-registered (unregistered branch arg)
	if len(store.GetBranches()) != 0 {
		t.Error("branch should NOT be auto-registered for unregistered branch arg")
	}
}

func TestExecuteWorktreeCreate_PathResolution(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	repoName := filepath.Base(repoDir)
	absBase := t.TempDir()

	template := absBase + "/{repository}/worktrees/{alias}"
	executeWorktreeCreate(&git, store, "path-test", "", "",
		template, repoDir, repoName)

	wt, found := store.GetWorktreeByAlias("path-test")
	if !found {
		t.Fatal("worktree should be registered")
	}

	expectedPath := filepath.Clean(absBase + "/" + repoName + "/worktrees/path-test")
	if wt.Path != expectedPath {
		t.Errorf("expected path %q, got %q", expectedPath, wt.Path)
	}
}
