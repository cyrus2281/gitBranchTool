package internal

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// initIntegrationRepo creates a temp git repo with an initial commit,
// plus a RepositoryBranches store in a separate temp dir.
func initIntegrationRepo(t *testing.T) (repoDir string, store *RepositoryBranches, git *Git) {
	t.Helper()

	repoDir = t.TempDir()
	storeDir := t.TempDir()

	// Init git repo
	cmds := [][]string{
		{"git", "init", repoDir},
		{"git", "-C", repoDir, "config", "user.email", "test@test.com"},
		{"git", "-C", repoDir, "config", "user.name", "Test"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v failed: %s\n%s", args, err, out)
		}
	}

	// Initial commit
	emptyFile := filepath.Join(repoDir, ".gitkeep")
	os.WriteFile(emptyFile, []byte(""), 0644)
	for _, args := range [][]string{
		{"git", "-C", repoDir, "add", ".gitkeep"},
		{"git", "-C", repoDir, "commit", "-m", "initial"},
	} {
		cmd := exec.Command(args[0], args[1:]...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v failed: %s\n%s", args, err, out)
		}
	}

	store = &RepositoryBranches{
		RepositoryName: filepath.Base(repoDir),
		StoreDirectory: storeDir,
	}

	git = &Git{}
	return repoDir, store, git
}

// chdir changes to a directory and returns a cleanup function.
func chdir(t *testing.T, dir string) {
	t.Helper()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	t.Cleanup(func() { os.Chdir(origDir) })
}

// --- Multi-step operation tests ---

func TestCreateBranchAndRegister(t *testing.T) {
	repoDir, store, git := initIntegrationRepo(t)
	chdir(t, repoDir)

	// Step 1: Create git branch
	err := git.SwitchToNewBranch("feature/auth")
	if err != nil {
		t.Fatalf("SwitchToNewBranch failed: %v", err)
	}

	// Step 2: Register in store
	branch := Branch{Name: "feature/auth", Alias: "auth", Note: "auth work"}
	store.AddBranch(branch)

	// Verify git state
	branches, err := git.GetBranches()
	if err != nil {
		t.Fatalf("GetBranches failed: %v", err)
	}
	gitHasBranch := false
	for _, b := range branches {
		if b == "feature/auth" {
			gitHasBranch = true
		}
	}
	if !gitHasBranch {
		t.Error("branch should exist in git")
	}

	// Verify store state
	if !store.BranchExists(Branch{Name: "feature/auth"}) {
		t.Error("branch should exist in store")
	}

	// Verify current branch
	g2 := Git{}
	current, _ := g2.GetCurrentBranch()
	if current != "feature/auth" {
		t.Errorf("expected current branch 'feature/auth', got %q", current)
	}
}

func TestCreateWorktreeAndRegister(t *testing.T) {
	repoDir, store, git := initIntegrationRepo(t)
	chdir(t, repoDir)

	// Step 1: Create branch + worktree
	wtPath := filepath.Join(repoDir, "worktrees", "my-wt")
	err := git.WorktreeAddNewBranch("feature/wt", wtPath)
	if err != nil {
		t.Fatalf("WorktreeAddNewBranch failed: %v", err)
	}

	// Step 2: Register branch
	store.AddBranch(Branch{Name: "feature/wt", Alias: "wt", Note: ""})

	// Step 3: Register worktree
	store.AddWorktree(Worktree{Alias: "my-wt", Path: wtPath, Note: ""})

	// Verify worktree directory exists
	if _, err := os.Stat(wtPath); os.IsNotExist(err) {
		t.Error("worktree directory should exist on disk")
	}

	// Verify git worktree list
	wtListOutput, err := git.WorktreeList()
	if err != nil {
		t.Fatalf("WorktreeList failed: %v", err)
	}
	wtMap := ParseWorktreeList(wtListOutput)

	// Resolve symlinks for path comparison (macOS /private/var)
	resolvedWtPath, _ := filepath.EvalSymlinks(wtPath)
	foundInGit := false
	for path, branch := range wtMap {
		resolvedPath, _ := filepath.EvalSymlinks(path)
		if resolvedPath == resolvedWtPath && branch == "feature/wt" {
			foundInGit = true
		}
	}
	if !foundInGit {
		t.Error("worktree should appear in git worktree list")
	}

	// Verify store state
	if !store.BranchExists(Branch{Name: "feature/wt"}) {
		t.Error("branch should be registered")
	}
	if !store.WorktreeAliasExists("my-wt") {
		t.Error("worktree should be registered")
	}
}

func TestDeleteBranchWithWorktree(t *testing.T) {
	repoDir, store, git := initIntegrationRepo(t)
	chdir(t, repoDir)

	// Setup: Create branch + worktree
	wtPath := filepath.Join(repoDir, "worktrees", "wt-del")
	git.WorktreeAddNewBranch("feature/del", wtPath)
	store.AddBranch(Branch{Name: "feature/del", Alias: "del", Note: ""})
	store.AddWorktree(Worktree{Alias: "wt-del", Path: wtPath, Note: ""})

	// Step 1: Delete worktree first (required before deleting branch)
	err := git.WorktreeRemove(wtPath, true)
	if err != nil {
		t.Fatalf("WorktreeRemove failed: %v", err)
	}

	// Remove worktree from store
	wt, found := store.GetWorktreeByAlias("wt-del")
	if found {
		store.RemoveWorktree(wt)
	}

	// Step 2: Delete branch
	err = git.DeleteBranch("feature/del", true)
	if err != nil {
		t.Fatalf("DeleteBranch failed: %v", err)
	}
	store.RemoveBranch(Branch{Name: "feature/del"})

	// Verify worktree dir is gone
	if _, err := os.Stat(wtPath); !os.IsNotExist(err) {
		t.Error("worktree directory should be removed")
	}

	// Verify git state: branch should be gone
	branches, _ := git.GetBranches()
	for _, b := range branches {
		if b == "feature/del" {
			t.Error("branch should be deleted from git")
		}
	}

	// Verify store state
	if store.BranchExists(Branch{Name: "feature/del"}) {
		t.Error("branch should be removed from store")
	}
	if store.WorktreeAliasExists("wt-del") {
		t.Error("worktree should be removed from store")
	}
}

func TestDeleteBranchFailsButWorktreeAlreadyDeleted(t *testing.T) {
	repoDir, store, git := initIntegrationRepo(t)
	chdir(t, repoDir)

	// Setup: Create branch + worktree
	wtPath := filepath.Join(repoDir, "worktrees", "wt-partial")
	git.WorktreeAddNewBranch("feature/partial", wtPath)
	store.AddBranch(Branch{Name: "feature/partial", Alias: "partial", Note: ""})
	store.AddWorktree(Worktree{Alias: "wt-partial", Path: wtPath, Note: ""})

	// Step 1: Delete worktree
	git.WorktreeRemove(wtPath, true)
	wt, _ := store.GetWorktreeByAlias("wt-partial")
	store.RemoveWorktree(wt)

	// Step 2: Try to delete worktree again (already deleted) — should not crash
	err := git.WorktreeRemove(wtPath, true)
	if err == nil {
		t.Log("worktree already deleted, error expected")
	}
	// This should handle gracefully — no panic

	// Verify worktree is gone from store
	if store.WorktreeAliasExists("wt-partial") {
		t.Error("worktree should still be removed from store")
	}

	// Branch should still be in store (since we didn't delete it yet)
	if !store.BranchExists(Branch{Name: "feature/partial"}) {
		t.Error("branch should still exist in store")
	}
}

func TestDeleteMultipleBranches(t *testing.T) {
	repoDir, store, git := initIntegrationRepo(t)
	chdir(t, repoDir)

	// Setup: Create 3 branches
	branchNames := []string{"feature/one", "feature/two", "feature/three"}
	for i, name := range branchNames {
		git.CreateNewBranch(name)
		store.AddBranch(Branch{Name: name, Alias: string(rune('a' + i)), Note: ""})
	}

	// Verify all exist
	if len(store.GetBranches()) != 3 {
		t.Fatalf("expected 3 branches, got %d", len(store.GetBranches()))
	}

	// Delete all in a loop (simulating multi-delete command)
	for _, name := range branchNames {
		branch, ok := store.GetBranchByName(name)
		if !ok {
			t.Fatalf("branch %q should exist before deletion", name)
		}

		err := git.DeleteBranch(name, true)
		if err != nil {
			t.Fatalf("DeleteBranch(%s) failed: %v", name, err)
		}
		store.RemoveBranch(branch)
	}

	// Verify all gone from store
	if len(store.GetBranches()) != 0 {
		t.Errorf("expected 0 branches after deletion, got %d", len(store.GetBranches()))
	}

	// Verify all gone from git
	branches, _ := git.GetBranches()
	branchSet := make(map[string]bool)
	for _, b := range branches {
		branchSet[b] = true
	}
	for _, name := range branchNames {
		if branchSet[name] {
			t.Errorf("branch %q should be deleted from git", name)
		}
	}
}

func TestPruneMultipleStaleWorktrees(t *testing.T) {
	repoDir, store, git := initIntegrationRepo(t)
	chdir(t, repoDir)

	// Setup: Create 3 worktrees
	wtPaths := make([]string, 3)
	for i := 0; i < 3; i++ {
		alias := string(rune('a' + i))
		wtPaths[i] = filepath.Join(repoDir, "worktrees", alias)
		branchName := "feature/" + alias
		git.WorktreeAddNewBranch(branchName, wtPaths[i])
		store.AddBranch(Branch{Name: branchName, Alias: alias, Note: ""})
		store.AddWorktree(Worktree{Alias: alias, Path: wtPaths[i], Note: ""})
	}

	if len(store.GetWorktrees()) != 3 {
		t.Fatalf("expected 3 worktrees, got %d", len(store.GetWorktrees()))
	}

	// Manually remove 2 worktree directories (simulating external deletion)
	// We need to use git worktree remove to properly clean up
	git.WorktreeRemove(wtPaths[0], true)
	git.WorktreeRemove(wtPaths[1], true)

	// Run prune equivalent: get active paths and remove stale entries
	git.WorktreePrune()

	wtListOutput, _ := git.WorktreeList()
	wtMap := ParseWorktreeList(wtListOutput)

	activePaths := make(map[string]bool)
	for path := range wtMap {
		// Resolve symlinks for macOS
		resolved, _ := filepath.EvalSymlinks(path)
		activePaths[resolved] = true
		activePaths[path] = true
	}

	// Collect stale entries (the fixed pattern from Bug 4)
	storedWorktrees := store.GetWorktrees()
	var toRemove []Worktree
	for _, wt := range storedWorktrees {
		resolvedPath, _ := filepath.EvalSymlinks(wt.Path)
		if !activePaths[wt.Path] && !activePaths[resolvedPath] {
			toRemove = append(toRemove, wt)
		}
	}
	for _, wt := range toRemove {
		store.RemoveWorktree(wt)
	}

	// KEY ASSERTION: Both stale worktrees should be pruned (Bug 4 fix)
	remaining := store.GetWorktrees()
	if len(remaining) != 1 {
		t.Errorf("expected 1 remaining worktree after pruning 2, got %d", len(remaining))
		for _, wt := range remaining {
			t.Logf("  remaining: alias=%s path=%s", wt.Alias, wt.Path)
		}
	}
}

func TestPruneNoStaleWorktrees(t *testing.T) {
	repoDir, store, git := initIntegrationRepo(t)
	chdir(t, repoDir)

	// Setup: Create a worktree (still valid)
	wtPath := filepath.Join(repoDir, "worktrees", "valid-wt")
	git.WorktreeAddNewBranch("feature/valid", wtPath)
	store.AddWorktree(Worktree{Alias: "valid-wt", Path: wtPath, Note: ""})

	// Get active paths
	wtListOutput, _ := git.WorktreeList()
	wtMap := ParseWorktreeList(wtListOutput)
	activePaths := make(map[string]bool)
	for path := range wtMap {
		resolved, _ := filepath.EvalSymlinks(path)
		activePaths[path] = true
		activePaths[resolved] = true
	}

	// Check for stale entries
	storedWorktrees := store.GetWorktrees()
	staleCount := 0
	for _, wt := range storedWorktrees {
		resolvedPath, _ := filepath.EvalSymlinks(wt.Path)
		if !activePaths[wt.Path] && !activePaths[resolvedPath] {
			staleCount++
		}
	}

	if staleCount != 0 {
		t.Errorf("expected 0 stale worktrees, got %d", staleCount)
	}
}

func TestWorktreeDeleteFailsGracefully(t *testing.T) {
	repoDir, store, git := initIntegrationRepo(t)
	chdir(t, repoDir)

	// Register a worktree that doesn't exist in git
	store.AddWorktree(Worktree{Alias: "ghost", Path: "/nonexistent/path", Note: ""})

	// Try to delete via git (should fail)
	err := git.WorktreeRemove("/nonexistent/path", false)
	if err == nil {
		t.Error("expected error when removing nonexistent worktree")
	}

	// Store should still have the entry (we didn't remove it since git failed)
	if !store.WorktreeAliasExists("ghost") {
		t.Error("worktree should still be in store after failed git removal")
	}

	// Now manually remove from store (like the command does on continue)
	// This simulates the command's behavior of continuing to next item
	if len(store.GetWorktrees()) != 1 {
		t.Errorf("expected 1 worktree in store, got %d", len(store.GetWorktrees()))
	}
}

func TestSwitchAndRegisterOnCheckout(t *testing.T) {
	repoDir, store, git := initIntegrationRepo(t)
	chdir(t, repoDir)

	// Create an unregistered branch
	git.CreateNewBranch("feature/unregistered")

	// Switch to it (simulating "g switch feature/unregistered myalias my note")
	err := git.SwitchBranch("feature/unregistered")
	if err != nil {
		t.Fatalf("SwitchBranch failed: %v", err)
	}

	// Register (simulating command behavior when alias is provided)
	alias := "unreg"
	if !store.AliasExists(alias) {
		store.AddBranch(Branch{
			Name:  "feature/unregistered",
			Alias: alias,
			Note:  "registered on checkout",
		})
	}

	// Verify current branch
	g2 := Git{}
	current, _ := g2.GetCurrentBranch()
	if current != "feature/unregistered" {
		t.Errorf("expected current branch 'feature/unregistered', got %q", current)
	}

	// Verify registered
	b, found := store.GetBranchByAlias("unreg")
	if !found {
		t.Error("branch should be registered after switch")
	}
	if b.Name != "feature/unregistered" {
		t.Errorf("expected name 'feature/unregistered', got '%s'", b.Name)
	}
}

func TestCreateBranchOnly_NoCheckout(t *testing.T) {
	repoDir, store, git := initIntegrationRepo(t)
	chdir(t, repoDir)

	// Get current branch before
	g1 := Git{}
	beforeBranch, _ := g1.GetCurrentBranch()

	// Create without checkout (simulating --only-create)
	err := git.CreateNewBranch("feature/no-checkout")
	if err != nil {
		t.Fatalf("CreateNewBranch failed: %v", err)
	}
	store.AddBranch(Branch{Name: "feature/no-checkout", Alias: "nc", Note: ""})

	// Verify we're still on original branch
	g2 := Git{}
	afterBranch, _ := g2.GetCurrentBranch()
	if afterBranch != beforeBranch {
		t.Errorf("should still be on %q, but switched to %q", beforeBranch, afterBranch)
	}

	// But branch should exist in git
	branches, _ := git.GetBranches()
	found := false
	for _, b := range branches {
		if b == "feature/no-checkout" {
			found = true
		}
	}
	if !found {
		t.Error("branch should exist in git even without checkout")
	}
}

func TestDeleteNonExistentBranch_StoreUnchanged(t *testing.T) {
	repoDir, store, git := initIntegrationRepo(t)
	chdir(t, repoDir)

	// Register a branch that doesn't exist in git anymore
	store.AddBranch(Branch{Name: "ghost-branch", Alias: "ghost", Note: ""})

	// Try to delete from git (should fail)
	err := git.DeleteBranch("ghost-branch", true)
	if err == nil {
		t.Error("expected error deleting non-existent git branch")
	}

	// Without --ignore-errors, store entry should remain
	if !store.BranchExists(Branch{Name: "ghost-branch"}) {
		t.Error("branch should remain in store when git delete fails")
	}

	// With --ignore-errors behavior: remove anyway
	store.RemoveBranch(Branch{Name: "ghost-branch"})
	if store.BranchExists(Branch{Name: "ghost-branch"}) {
		t.Error("branch should be removed from store with ignore-errors")
	}
}

func TestWorktreeCreation_PathResolution(t *testing.T) {
	repoDir, _, git := initIntegrationRepo(t)
	chdir(t, repoDir)

	// Test ResolveWorktreePath with real repo values
	template := DEFAULT_WORKTREE_PATH
	repoName := filepath.Base(repoDir)

	resolvedPath := ResolveWorktreePath(template, repoDir, repoName, "test-alias", "feature/test")
	expectedPath := filepath.Join(repoDir, "worktrees", "test-alias")

	if resolvedPath != expectedPath {
		t.Errorf("expected %q, got %q", expectedPath, resolvedPath)
	}

	// Actually create the worktree at the resolved path
	err := git.WorktreeAddNewBranch("feature/test", resolvedPath)
	if err != nil {
		t.Fatalf("WorktreeAddNewBranch failed: %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(resolvedPath); os.IsNotExist(err) {
		t.Error("worktree directory should exist at resolved path")
	}
}
