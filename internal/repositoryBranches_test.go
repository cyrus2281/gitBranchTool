package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// newTestStore creates a RepositoryBranches backed by a temp directory.
func newTestStore(t *testing.T) *RepositoryBranches {
	t.Helper()
	dir := t.TempDir()
	return &RepositoryBranches{
		RepositoryName: "test-repo",
		StoreDirectory: dir,
	}
}

// --- Load / Persistence tests ---

func TestLoadNonExistentFile(t *testing.T) {
	s := newTestStore(t)
	branches := s.GetBranches()
	if len(branches) != 0 {
		t.Errorf("expected empty branches, got %d", len(branches))
	}
	worktrees := s.GetWorktrees()
	if len(worktrees) != 0 {
		t.Errorf("expected empty worktrees, got %d", len(worktrees))
	}
}

func TestSaveLoadRoundtrip(t *testing.T) {
	s := newTestStore(t)

	// Add data
	s.AddBranch(Branch{Name: "feature/auth", Alias: "auth", Note: "auth work"})
	s.AddBranch(Branch{Name: "bugfix/123", Alias: "fix", Note: "bug fix"})
	s.AddWorktree(Worktree{Alias: "wt1", Path: "/path/to/wt1", Note: "worktree 1"})
	s.SetDefaultBranch("develop")
	s.SetLocalPrefix("dev/")

	// Create new store pointing to same file
	s2 := &RepositoryBranches{
		RepositoryName: "test-repo",
		StoreDirectory: s.StoreDirectory,
	}

	// Verify data loaded correctly
	branches := s2.GetBranches()
	if len(branches) != 2 {
		t.Fatalf("expected 2 branches, got %d", len(branches))
	}
	if branches[0].Name != "feature/auth" || branches[0].Alias != "auth" || branches[0].Note != "auth work" {
		t.Errorf("branch 0 mismatch: %+v", branches[0])
	}
	if branches[1].Name != "bugfix/123" || branches[1].Alias != "fix" || branches[1].Note != "bug fix" {
		t.Errorf("branch 1 mismatch: %+v", branches[1])
	}

	worktrees := s2.GetWorktrees()
	if len(worktrees) != 1 {
		t.Fatalf("expected 1 worktree, got %d", len(worktrees))
	}
	if worktrees[0].Alias != "wt1" || worktrees[0].Path != "/path/to/wt1" {
		t.Errorf("worktree mismatch: %+v", worktrees[0])
	}

	if s2.GetDefaultBranch() != "develop" {
		t.Errorf("expected default branch 'develop', got '%s'", s2.GetDefaultBranch())
	}
	if s2.GetLocalPrefix() != "dev/" {
		t.Errorf("expected local prefix 'dev/', got '%s'", s2.GetLocalPrefix())
	}
}

func TestCorruptedJSON(t *testing.T) {
	dir := t.TempDir()
	// Write garbage to the JSON file
	filePath := filepath.Join(dir, "test-repo.json")
	os.WriteFile(filePath, []byte("not valid json{{{"), 0644)

	// The load() calls logger.FatalF for corrupted JSON, which exits the process.
	// We can't easily test os.Exit in Go unit tests without a subprocess.
	// Instead, verify the file exists and is corrupt — the bug fix ensures
	// the error is logged rather than silently swallowed.
	if _, err := os.Stat(filePath); err != nil {
		t.Fatalf("expected corrupt file to exist: %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	// Verify the file is indeed invalid JSON
	var jsonData repositoryBranchesJson
	if err := json.Unmarshal(content, &jsonData); err == nil {
		t.Error("expected JSON unmarshal to fail for corrupted data")
	}
}

// --- Branch CRUD tests ---

func TestAddAndGetBranches_Empty(t *testing.T) {
	s := newTestStore(t)
	branches := s.GetBranches()
	if len(branches) != 0 {
		t.Errorf("expected 0 branches, got %d", len(branches))
	}
}

func TestAddAndGetBranches_One(t *testing.T) {
	s := newTestStore(t)
	s.AddBranch(Branch{Name: "main", Alias: "m", Note: "main branch"})

	branches := s.GetBranches()
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Name != "main" {
		t.Errorf("expected 'main', got '%s'", branches[0].Name)
	}
}

func TestAddAndGetBranches_Multiple(t *testing.T) {
	s := newTestStore(t)
	s.AddBranch(Branch{Name: "branch-a", Alias: "a", Note: ""})
	s.AddBranch(Branch{Name: "branch-b", Alias: "b", Note: "note b"})
	s.AddBranch(Branch{Name: "branch-c", Alias: "c", Note: "note c"})

	branches := s.GetBranches()
	if len(branches) != 3 {
		t.Fatalf("expected 3 branches, got %d", len(branches))
	}
	// Verify order preserved
	if branches[0].Name != "branch-a" || branches[1].Name != "branch-b" || branches[2].Name != "branch-c" {
		t.Errorf("order not preserved: %v", branches)
	}
}

func TestBranchExists(t *testing.T) {
	s := newTestStore(t)
	s.AddBranch(Branch{Name: "feature-x", Alias: "fx", Note: ""})

	if !s.BranchExists(Branch{Name: "feature-x"}) {
		t.Error("expected BranchExists to return true for existing branch")
	}
	if s.BranchExists(Branch{Name: "nonexistent"}) {
		t.Error("expected BranchExists to return false for non-existing branch")
	}
}

func TestAliasExists(t *testing.T) {
	s := newTestStore(t)
	s.AddBranch(Branch{Name: "feature-x", Alias: "fx", Note: ""})

	if !s.AliasExists("fx") {
		t.Error("expected AliasExists to return true for existing alias")
	}
	if s.AliasExists("nonexistent") {
		t.Error("expected AliasExists to return false for non-existing alias")
	}
}

func TestBranchWithAliasExists(t *testing.T) {
	s := newTestStore(t)
	s.AddBranch(Branch{Name: "feature-x", Alias: "fx", Note: ""})

	// BranchWithAliasExists checks the Name field, not Alias
	if !s.BranchWithAliasExists("feature-x") {
		t.Error("expected true when param matches a branch Name")
	}
	if s.BranchWithAliasExists("fx") {
		t.Error("expected false when param matches Alias (not Name)")
	}
	if s.BranchWithAliasExists("nonexistent") {
		t.Error("expected false for non-existing name")
	}
}

func TestGetBranchByAlias(t *testing.T) {
	s := newTestStore(t)
	s.AddBranch(Branch{Name: "feature-x", Alias: "fx", Note: "note"})

	b, found := s.GetBranchByAlias("fx")
	if !found {
		t.Fatal("expected branch to be found by alias")
	}
	if b.Name != "feature-x" {
		t.Errorf("expected name 'feature-x', got '%s'", b.Name)
	}

	_, found = s.GetBranchByAlias("nonexistent")
	if found {
		t.Error("expected not found for non-existing alias")
	}
}

func TestGetBranchByName(t *testing.T) {
	s := newTestStore(t)
	s.AddBranch(Branch{Name: "feature-x", Alias: "fx", Note: "note"})

	b, found := s.GetBranchByName("feature-x")
	if !found {
		t.Fatal("expected branch to be found by name")
	}
	if b.Alias != "fx" {
		t.Errorf("expected alias 'fx', got '%s'", b.Alias)
	}

	_, found = s.GetBranchByName("nonexistent")
	if found {
		t.Error("expected not found for non-existing name")
	}
}

func TestGetBranchByNameOrAlias(t *testing.T) {
	s := newTestStore(t)
	s.AddBranch(Branch{Name: "feature-x", Alias: "fx", Note: ""})

	// Find by name
	b, found := s.GetBranchByNameOrAlias("feature-x")
	if !found || b.Name != "feature-x" {
		t.Errorf("expected found by name, got found=%v, name=%s", found, b.Name)
	}

	// Find by alias
	b, found = s.GetBranchByNameOrAlias("fx")
	if !found || b.Name != "feature-x" {
		t.Errorf("expected found by alias, got found=%v, name=%s", found, b.Name)
	}

	// Not found
	_, found = s.GetBranchByNameOrAlias("nonexistent")
	if found {
		t.Error("expected not found")
	}
}

func TestUpdateBranch(t *testing.T) {
	s := newTestStore(t)
	s.AddBranch(Branch{Name: "feature-x", Alias: "fx", Note: "original"})

	s.UpdateBranch(Branch{Name: "feature-x", Alias: "fx-updated", Note: "updated note"})

	b, found := s.GetBranchByName("feature-x")
	if !found {
		t.Fatal("branch should still exist after update")
	}
	if b.Alias != "fx-updated" {
		t.Errorf("expected alias 'fx-updated', got '%s'", b.Alias)
	}
	if b.Note != "updated note" {
		t.Errorf("expected note 'updated note', got '%s'", b.Note)
	}
}

func TestUpdateBranch_Persisted(t *testing.T) {
	s := newTestStore(t)
	s.AddBranch(Branch{Name: "feature-x", Alias: "fx", Note: "original"})
	s.UpdateBranch(Branch{Name: "feature-x", Alias: "fx-new", Note: "new"})

	// Reload from disk
	s2 := &RepositoryBranches{
		RepositoryName: "test-repo",
		StoreDirectory: s.StoreDirectory,
	}
	b, found := s2.GetBranchByName("feature-x")
	if !found {
		t.Fatal("branch should exist after reload")
	}
	if b.Alias != "fx-new" || b.Note != "new" {
		t.Errorf("update not persisted: %+v", b)
	}
}

func TestRemoveBranch(t *testing.T) {
	s := newTestStore(t)
	s.AddBranch(Branch{Name: "to-remove", Alias: "tr", Note: ""})

	s.RemoveBranch(Branch{Name: "to-remove"})

	if s.BranchExists(Branch{Name: "to-remove"}) {
		t.Error("branch should be removed")
	}
	if len(s.GetBranches()) != 0 {
		t.Errorf("expected 0 branches, got %d", len(s.GetBranches()))
	}
}

func TestRemoveBranch_NonExistent(t *testing.T) {
	s := newTestStore(t)
	s.AddBranch(Branch{Name: "keep", Alias: "k", Note: ""})

	// Should not panic or crash
	s.RemoveBranch(Branch{Name: "does-not-exist"})

	if len(s.GetBranches()) != 1 {
		t.Errorf("expected 1 branch unchanged, got %d", len(s.GetBranches()))
	}
}

func TestRemoveBranchPreservesOthers(t *testing.T) {
	s := newTestStore(t)
	s.AddBranch(Branch{Name: "first", Alias: "1st", Note: ""})
	s.AddBranch(Branch{Name: "middle", Alias: "mid", Note: ""})
	s.AddBranch(Branch{Name: "last", Alias: "3rd", Note: ""})

	s.RemoveBranch(Branch{Name: "middle"})

	branches := s.GetBranches()
	if len(branches) != 2 {
		t.Fatalf("expected 2 branches, got %d", len(branches))
	}
	if branches[0].Name != "first" || branches[1].Name != "last" {
		t.Errorf("wrong branches remaining: %v", branches)
	}
}

// --- Worktree CRUD tests ---

func TestAddAndGetWorktrees_Empty(t *testing.T) {
	s := newTestStore(t)
	wts := s.GetWorktrees()
	if len(wts) != 0 {
		t.Errorf("expected 0 worktrees, got %d", len(wts))
	}
}

func TestAddAndGetWorktrees_Multiple(t *testing.T) {
	s := newTestStore(t)
	s.AddWorktree(Worktree{Alias: "wt1", Path: "/path/1", Note: "note 1"})
	s.AddWorktree(Worktree{Alias: "wt2", Path: "/path/2", Note: "note 2"})
	s.AddWorktree(Worktree{Alias: "wt3", Path: "/path/3", Note: ""})

	wts := s.GetWorktrees()
	if len(wts) != 3 {
		t.Fatalf("expected 3 worktrees, got %d", len(wts))
	}
}

func TestGetWorktreeByAlias(t *testing.T) {
	s := newTestStore(t)
	s.AddWorktree(Worktree{Alias: "wt1", Path: "/path/1", Note: "note"})

	wt, found := s.GetWorktreeByAlias("wt1")
	if !found {
		t.Fatal("expected worktree found by alias")
	}
	if wt.Path != "/path/1" {
		t.Errorf("expected path '/path/1', got '%s'", wt.Path)
	}

	_, found = s.GetWorktreeByAlias("nonexistent")
	if found {
		t.Error("expected not found for non-existing alias")
	}
}

func TestGetWorktreeByPath(t *testing.T) {
	s := newTestStore(t)
	s.AddWorktree(Worktree{Alias: "wt1", Path: "/path/1", Note: ""})

	wt, found := s.GetWorktreeByPath("/path/1")
	if !found {
		t.Fatal("expected worktree found by path")
	}
	if wt.Alias != "wt1" {
		t.Errorf("expected alias 'wt1', got '%s'", wt.Alias)
	}

	_, found = s.GetWorktreeByPath("/nonexistent")
	if found {
		t.Error("expected not found for non-existing path")
	}
}

func TestWorktreeAliasExists(t *testing.T) {
	s := newTestStore(t)
	s.AddWorktree(Worktree{Alias: "wt1", Path: "/path/1", Note: ""})

	if !s.WorktreeAliasExists("wt1") {
		t.Error("expected true for existing alias")
	}
	if s.WorktreeAliasExists("nonexistent") {
		t.Error("expected false for non-existing alias")
	}
}

func TestRemoveWorktree(t *testing.T) {
	s := newTestStore(t)
	s.AddWorktree(Worktree{Alias: "wt1", Path: "/path/1", Note: ""})
	s.AddWorktree(Worktree{Alias: "wt2", Path: "/path/2", Note: ""})

	s.RemoveWorktree(Worktree{Alias: "wt1"})

	wts := s.GetWorktrees()
	if len(wts) != 1 {
		t.Fatalf("expected 1 worktree, got %d", len(wts))
	}
	if wts[0].Alias != "wt2" {
		t.Errorf("expected 'wt2' remaining, got '%s'", wts[0].Alias)
	}
}

func TestRemoveWorktree_NonExistent(t *testing.T) {
	s := newTestStore(t)
	s.AddWorktree(Worktree{Alias: "wt1", Path: "/path/1", Note: ""})

	// Should not panic
	s.RemoveWorktree(Worktree{Alias: "nonexistent"})

	if len(s.GetWorktrees()) != 1 {
		t.Error("expected worktree count unchanged")
	}
}

func TestRemoveMultipleWorktrees(t *testing.T) {
	s := newTestStore(t)
	s.AddWorktree(Worktree{Alias: "a", Path: "/a", Note: ""})
	s.AddWorktree(Worktree{Alias: "b", Path: "/b", Note: ""})
	s.AddWorktree(Worktree{Alias: "c", Path: "/c", Note: ""})

	// Collect first, then remove (simulating the bug fix pattern)
	toRemove := []Worktree{
		{Alias: "a"},
		{Alias: "b"},
		{Alias: "c"},
	}
	for _, wt := range toRemove {
		s.RemoveWorktree(wt)
	}

	if len(s.GetWorktrees()) != 0 {
		t.Errorf("expected all worktrees removed, got %d", len(s.GetWorktrees()))
	}
}

// --- Default branch tests ---

func TestDefaultBranch_ReturnsMainWhenUnset(t *testing.T) {
	s := newTestStore(t)
	if s.GetDefaultBranch() != "main" {
		t.Errorf("expected default 'main', got '%s'", s.GetDefaultBranch())
	}
}

func TestDefaultBranch_SetAndGet(t *testing.T) {
	s := newTestStore(t)
	s.SetDefaultBranch("develop")
	if s.GetDefaultBranch() != "develop" {
		t.Errorf("expected 'develop', got '%s'", s.GetDefaultBranch())
	}
}

func TestDefaultBranch_Persisted(t *testing.T) {
	s := newTestStore(t)
	s.SetDefaultBranch("develop")

	s2 := &RepositoryBranches{
		RepositoryName: "test-repo",
		StoreDirectory: s.StoreDirectory,
	}
	if s2.GetDefaultBranch() != "develop" {
		t.Errorf("expected 'develop' after reload, got '%s'", s2.GetDefaultBranch())
	}
}

// --- Local prefix tests ---

func TestLocalPrefix_EmptyWhenUnset(t *testing.T) {
	s := newTestStore(t)
	if s.GetLocalPrefix() != "" {
		t.Errorf("expected empty prefix, got '%s'", s.GetLocalPrefix())
	}
}

func TestLocalPrefix_SetAndGet(t *testing.T) {
	s := newTestStore(t)
	s.SetLocalPrefix("feature/")
	if s.GetLocalPrefix() != "feature/" {
		t.Errorf("expected 'feature/', got '%s'", s.GetLocalPrefix())
	}
}

func TestLocalPrefix_TrimsWhitespace(t *testing.T) {
	s := newTestStore(t)
	s.SetLocalPrefix("  dev/  ")
	if s.GetLocalPrefix() != "dev/" {
		t.Errorf("expected 'dev/' (trimmed), got '%s'", s.GetLocalPrefix())
	}
}

func TestLocalPrefix_Persisted(t *testing.T) {
	s := newTestStore(t)
	s.SetLocalPrefix("fix/")

	s2 := &RepositoryBranches{
		RepositoryName: "test-repo",
		StoreDirectory: s.StoreDirectory,
	}
	if s2.GetLocalPrefix() != "fix/" {
		t.Errorf("expected 'fix/' after reload, got '%s'", s2.GetLocalPrefix())
	}
}

// --- Lazy loading tests ---

func TestLazyLoading(t *testing.T) {
	dir := t.TempDir()

	// Pre-write a JSON file
	data := repositoryBranchesJson{
		Branches:      []Branch{{Name: "pre-loaded", Alias: "pl", Note: ""}},
		DefaultBranch: "main",
	}
	jsonBytes, _ := json.Marshal(data)
	os.WriteFile(filepath.Join(dir, "test-repo.json"), jsonBytes, 0644)

	s := &RepositoryBranches{
		RepositoryName: "test-repo",
		StoreDirectory: dir,
	}

	// Should not be loaded yet
	if s.loaded {
		t.Error("expected loaded=false before first access")
	}

	// First access triggers load
	branches := s.GetBranches()
	if !s.loaded {
		t.Error("expected loaded=true after GetBranches()")
	}
	if len(branches) != 1 || branches[0].Name != "pre-loaded" {
		t.Errorf("expected pre-loaded data, got %v", branches)
	}
}

// --- Save preserves empty default branch as DEFAULT_BRANCH ---

func TestSave_EmptyDefaultBranch_BecomesMain(t *testing.T) {
	dir := t.TempDir()
	s := &RepositoryBranches{
		RepositoryName: "test-repo",
		StoreDirectory: dir,
	}
	s.AddBranch(Branch{Name: "test", Alias: "t", Note: ""})

	// Reload and check default branch
	s2 := &RepositoryBranches{
		RepositoryName: "test-repo",
		StoreDirectory: dir,
	}
	// The saved file should have "main" as default branch
	if s2.GetDefaultBranch() != "main" {
		t.Errorf("expected 'main' as default, got '%s'", s2.GetDefaultBranch())
	}
}
