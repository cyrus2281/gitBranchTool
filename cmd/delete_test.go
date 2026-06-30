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

func TestResolveDeleteTargets_NoRegexReturnsArgsUnchanged(t *testing.T) {
	store := newTestStore(t)
	store.AddBranch(internal.Branch{Name: "feature/a", Alias: "a"})

	args := []string{"feature/a", "b", "nonexistent"}
	targets, err := resolveDeleteTargets(store, nil, args, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(targets) != len(args) {
		t.Fatalf("expected args returned unchanged, got %v", targets)
	}
	for i, v := range args {
		if targets[i] != v {
			t.Errorf("expected targets[%d]=%q, got %q", i, v, targets[i])
		}
	}
}

func TestResolveDeleteTargets_MatchesByName(t *testing.T) {
	store := newTestStore(t)
	store.AddBranch(internal.Branch{Name: "claude/one", Alias: "c1"})
	store.AddBranch(internal.Branch{Name: "claude/two", Alias: "c2"})
	store.AddBranch(internal.Branch{Name: "feature/keep", Alias: "keep"})

	targets, err := resolveDeleteTargets(store, nil, []string{"claude.*"}, true, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(targets) != 2 {
		t.Fatalf("expected 2 matches, got %d (%v)", len(targets), targets)
	}
	for _, name := range targets {
		if name == "feature/keep" {
			t.Error("non-matching branch should not be a target")
		}
	}
}

func TestResolveDeleteTargets_MatchesByAlias(t *testing.T) {
	store := newTestStore(t)
	store.AddBranch(internal.Branch{Name: "feature/login", Alias: "rel-1/beta"})
	store.AddBranch(internal.Branch{Name: "feature/logout", Alias: "rel-9/beta"})

	targets, err := resolveDeleteTargets(store, nil, []string{"rel-[1-3]/beta"}, true, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(targets) != 1 || targets[0] != "feature/login" {
		t.Fatalf("expected only feature/login matched by alias, got %v", targets)
	}
}

func TestResolveDeleteTargets_MultiplePatternsDeduped(t *testing.T) {
	store := newTestStore(t)
	store.AddBranch(internal.Branch{Name: "cyrus/jira-1", Alias: "j1"})
	store.AddBranch(internal.Branch{Name: "cyrus/jira-2", Alias: "j2"})
	store.AddBranch(internal.Branch{Name: "rel-2/beta", Alias: "rb"})

	// Two patterns; "cyrus/.*" also overlaps nothing with the second, but the
	// branch order must be preserved and no branch may appear twice.
	targets, err := resolveDeleteTargets(store, nil, []string{"cyrus/.*", "rel-[1-3]/beta"}, true, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(targets) != 3 {
		t.Fatalf("expected 3 matches, got %d (%v)", len(targets), targets)
	}
	expected := []string{"cyrus/jira-1", "cyrus/jira-2", "rel-2/beta"}
	for i, name := range expected {
		if targets[i] != name {
			t.Errorf("expected targets[%d]=%q, got %q", i, name, targets[i])
		}
	}
}

func TestResolveDeleteTargets_NoMatches(t *testing.T) {
	store := newTestStore(t)
	store.AddBranch(internal.Branch{Name: "feature/a", Alias: "a"})

	targets, err := resolveDeleteTargets(store, nil, []string{"nomatch.*"}, true, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(targets) != 0 {
		t.Errorf("expected no matches, got %v", targets)
	}
}

func TestResolveDeleteTargets_InvalidRegex(t *testing.T) {
	store := newTestStore(t)
	store.AddBranch(internal.Branch{Name: "feature/a", Alias: "a"})

	_, err := resolveDeleteTargets(store, nil, []string{"["}, true, false)
	if err == nil {
		t.Error("expected an error for an invalid regular expression")
	}
}

func TestResolveDeleteTargets_AllRegexMatchesUnregistered(t *testing.T) {
	store := newTestStore(t)
	store.AddBranch(internal.Branch{Name: "claude/one", Alias: "c1"})

	gitBranches := []string{"claude/one", "claude/two", "main"}
	targets, err := resolveDeleteTargets(store, gitBranches, []string{"claude/.*"}, true, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Registered "claude/one" comes first, then the unregistered "claude/two".
	expected := []string{"claude/one", "claude/two"}
	if len(targets) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, targets)
	}
	for i, name := range expected {
		if targets[i] != name {
			t.Errorf("expected targets[%d]=%q, got %q", i, name, targets[i])
		}
	}
}

func TestResolveDeleteTargets_AllFalseIgnoresUnregistered(t *testing.T) {
	store := newTestStore(t)
	store.AddBranch(internal.Branch{Name: "claude/one", Alias: "c1"})

	gitBranches := []string{"claude/one", "claude/two"}
	// all=false → the unregistered "claude/two" must not be matched.
	targets, err := resolveDeleteTargets(store, gitBranches, []string{"claude/.*"}, true, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(targets) != 1 || targets[0] != "claude/one" {
		t.Fatalf("expected only registered claude/one, got %v", targets)
	}
}

func TestExecuteDeleteBranch_AllDeletesUnregistered(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	// Branch exists in git but is not registered with g.
	git.CreateNewBranch("feature/unregistered")

	executeDeleteBranch(&git, store, "feature/unregistered", deleteOpts{Force: true, All: true})

	branches, _ := git.GetBranches()
	for _, b := range branches {
		if b == "feature/unregistered" {
			t.Error("unregistered branch should be deleted from git with --all")
		}
	}
}

func TestExecuteDeleteBranch_UnregisteredSkippedWithoutAll(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	git.CreateNewBranch("feature/unregistered")

	// Without --all an unregistered branch is reported as not found and kept.
	executeDeleteBranch(&git, store, "feature/unregistered", deleteOpts{Force: true})

	branches, _ := git.GetBranches()
	found := false
	for _, b := range branches {
		if b == "feature/unregistered" {
			found = true
		}
	}
	if !found {
		t.Error("unregistered branch should remain when --all is not set")
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
