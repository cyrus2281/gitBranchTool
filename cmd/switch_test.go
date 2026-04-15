package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cyrus2281/gitBranchTool/internal"
)

func TestExecuteSwitch_RegisteredByAlias(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	git.CreateNewBranch("feature/registered")
	store.AddBranch(internal.Branch{Name: "feature/registered", Alias: "reg", Note: ""})

	executeSwitch(&git, store, "reg", "", "", switchOpts{})

	g2 := internal.Git{}
	current, _ := g2.GetCurrentBranch()
	if current != "feature/registered" {
		t.Errorf("expected checkout to 'feature/registered', got %q", current)
	}
}

func TestExecuteSwitch_UnregisteredWithAlias(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	git.CreateNewBranch("feature/unreg")

	executeSwitch(&git, store, "feature/unreg", "ur", "my note ", switchOpts{})

	// Should be registered now
	b, found := store.GetBranchByAlias("ur")
	if !found {
		t.Error("branch should be registered after switch with alias")
	}
	if b.Name != "feature/unreg" {
		t.Errorf("expected name 'feature/unreg', got %q", b.Name)
	}
	if b.Note != "my note " {
		t.Errorf("expected note 'my note ', got %q", b.Note)
	}
}

func TestExecuteSwitch_UnregisteredWithoutAlias(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	git.CreateNewBranch("feature/no-alias")

	executeSwitch(&git, store, "feature/no-alias", "", "", switchOpts{})

	// Should NOT be registered
	if len(store.GetBranches()) != 0 {
		t.Error("branch should not be registered when no alias provided")
	}
}

func TestExecuteSwitch_AliasConflict(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	git.CreateNewBranch("feature/conflict")
	store.AddBranch(internal.Branch{Name: "other", Alias: "taken", Note: ""})

	// Switch should succeed but NOT register (alias conflict warning)
	executeSwitch(&git, store, "feature/conflict", "taken", "", switchOpts{})

	// Should have switched
	g2 := internal.Git{}
	current, _ := g2.GetCurrentBranch()
	if current != "feature/conflict" {
		t.Errorf("should have switched to 'feature/conflict', got %q", current)
	}

	// Should only have the original branch registered (conflict prevented new registration)
	if len(store.GetBranches()) != 1 {
		t.Errorf("expected 1 branch (original only), got %d", len(store.GetBranches()))
	}
}

func TestExecuteSwitch_WorktreeMode_ExistingWorktree(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	// Create branch + worktree first
	wtPath := filepath.Join(repoDir, "worktrees", "existing")
	git.WorktreeAddNewBranch("feature/existing-wt", wtPath)

	executeSwitch(&git, store, "feature/existing-wt", "", "", switchOpts{
		UseWorktree:          true,
		WorktreePathTemplate: "./worktrees/{alias}",
		RepoPath:             repoDir,
		RepoName:             filepath.Base(repoDir),
	})

	// Should NOT create a new worktree — just report existing one
}

func TestExecuteSwitch_WorktreeMode_CreateNew(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	git.CreateNewBranch("feature/new-wt")

	executeSwitch(&git, store, "feature/new-wt", "nw", "", switchOpts{
		UseWorktree:          true,
		WorktreeAlias:        "new-wt-alias",
		WorktreePathTemplate: "./worktrees/{alias}",
		RepoPath:             repoDir,
		RepoName:             filepath.Base(repoDir),
	})

	// Worktree should be registered
	if !store.WorktreeAliasExists("new-wt-alias") {
		t.Error("worktree should be registered")
	}

	// Worktree directory should exist
	wtPath := filepath.Join(repoDir, "worktrees", "new-wt-alias")
	if _, err := os.Stat(wtPath); os.IsNotExist(err) {
		t.Error("worktree directory should exist")
	}
}
