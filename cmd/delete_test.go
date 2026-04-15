package cmd

import (
	"os/exec"
	"testing"

	"github.com/cyrus2281/gitBranchTool/internal"
)

func TestExecuteDeleteBranch_ByName(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	git.CreateNewBranch("feature/del")
	store.AddBranch(internal.Branch{Name: "feature/del", Alias: "del", Note: ""})

	executeDeleteBranch(&git, store, "feature/del", deleteOpts{Force: true})

	// Branch should be removed from store
	if store.BranchExists(internal.Branch{Name: "feature/del"}) {
		t.Error("branch should be removed from store")
	}

	// Branch should be removed from git
	branches, _ := git.GetBranches()
	for _, b := range branches {
		if b == "feature/del" {
			t.Error("branch should be deleted from git")
		}
	}
}

func TestExecuteDeleteBranch_ByAlias(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	git.CreateNewBranch("feature/aliased")
	store.AddBranch(internal.Branch{Name: "feature/aliased", Alias: "al", Note: ""})

	executeDeleteBranch(&git, store, "al", deleteOpts{Force: true})

	if store.BranchExists(internal.Branch{Name: "feature/aliased"}) {
		t.Error("branch should be removed when deleted by alias")
	}
}

func TestExecuteDeleteBranch_NotFound(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	// Should not panic when branch not found
	executeDeleteBranch(&git, store, "nonexistent", deleteOpts{Force: true})
}

func TestExecuteDeleteBranch_GitFailsWithoutIgnoreErrors(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	// Register a branch that doesn't exist in git
	store.AddBranch(internal.Branch{Name: "ghost-branch", Alias: "ghost", Note: ""})

	executeDeleteBranch(&git, store, "ghost-branch", deleteOpts{Force: true})

	// Without IgnoreErrors, store entry should remain
	if !store.BranchExists(internal.Branch{Name: "ghost-branch"}) {
		t.Error("branch should remain in store when git delete fails without ignore-errors")
	}
}

func TestExecuteDeleteBranch_GitFailsWithIgnoreErrors(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	// Register a branch that doesn't exist in git
	store.AddBranch(internal.Branch{Name: "ghost-branch", Alias: "ghost", Note: ""})

	executeDeleteBranch(&git, store, "ghost-branch", deleteOpts{Force: true, IgnoreErrors: true})

	// With IgnoreErrors, store entry should be removed anyway
	if store.BranchExists(internal.Branch{Name: "ghost-branch"}) {
		t.Error("branch should be removed from store with ignore-errors even when git fails")
	}
}

func TestExecuteDeleteBranch_RemoteOnly(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	git.CreateNewBranch("feature/remote-only")
	store.AddBranch(internal.Branch{Name: "feature/remote-only", Alias: "ro", Note: ""})

	executeDeleteBranch(&git, store, "feature/remote-only", deleteOpts{RemoteOnly: true})

	// Local branch and store entry should remain
	if !store.BranchExists(internal.Branch{Name: "feature/remote-only"}) {
		t.Error("store entry should remain with remote-only")
	}
	branches, _ := git.GetBranches()
	found := false
	for _, b := range branches {
		if b == "feature/remote-only" {
			found = true
		}
	}
	if !found {
		t.Error("local git branch should remain with remote-only")
	}
}

func TestExecuteDeleteBranch_WithWorktree(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	// Create branch + worktree
	wtPath := repoDir + "/worktrees/wt-del"
	git.WorktreeAddNewBranch("feature/wt-del", wtPath)
	store.AddBranch(internal.Branch{Name: "feature/wt-del", Alias: "wtd", Note: ""})

	// Use the git-reported path for the store (resolves symlinks on macOS)
	wtOutput, _ := git.WorktreeList()
	wtMap := internal.ParseWorktreeList(wtOutput)
	resolvedWtPath := internal.GetWorktreePathForBranch(wtMap, "feature/wt-del")
	store.AddWorktree(internal.Worktree{Alias: "wt-del", Path: resolvedWtPath, Note: ""})

	executeDeleteBranch(&git, store, "feature/wt-del", deleteOpts{
		Force:                true,
		ShouldDeleteWorktree: true,
		WorktreeMap:          wtMap,
	})

	// Both should be cleaned up
	if store.BranchExists(internal.Branch{Name: "feature/wt-del"}) {
		t.Error("branch should be removed from store")
	}
	if store.WorktreeAliasExists("wt-del") {
		t.Error("worktree should be removed from store")
	}
}

func TestExecuteDeleteBranch_MultipleInSequence(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	names := []string{"feature/a", "feature/b", "feature/c"}
	for i, name := range names {
		git.CreateNewBranch(name)
		store.AddBranch(internal.Branch{Name: name, Alias: string(rune('a' + i)), Note: ""})
	}

	for _, name := range names {
		executeDeleteBranch(&git, store, name, deleteOpts{Force: true})
	}

	if len(store.GetBranches()) != 0 {
		t.Errorf("expected 0 branches after deleting all, got %d", len(store.GetBranches()))
	}
}

func TestExecuteDeleteBranch_SafeDelete(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	// Create an unmerged branch with a commit
	git.SwitchToNewBranch("unmerged-feature")
	exec.Command("git", "-C", repoDir, "commit", "--allow-empty", "-m", "unmerged commit").Run()
	git2 := internal.Git{}
	defaultBranch, _ := git2.GetCurrentBranch()
	exec.Command("git", "-C", repoDir, "checkout", "-").Run()

	store.AddBranch(internal.Branch{Name: "unmerged-feature", Alias: "uf", Note: ""})

	// Safe delete (Force=false) should fail for unmerged branch
	executeDeleteBranch(&internal.Git{}, store, "unmerged-feature", deleteOpts{Force: false})

	// Branch should remain because safe delete failed (unmerged)
	if !store.BranchExists(internal.Branch{Name: "unmerged-feature"}) {
		t.Error("branch should remain after failed safe delete")
	}
	_ = defaultBranch
}
