package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/cyrus2281/gitBranchTool/internal"
)

func initTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, args := range [][]string{
		{"git", "init", dir},
		{"git", "-C", dir, "config", "user.email", "test@test.com"},
		{"git", "-C", dir, "config", "user.name", "Test"},
	} {
		cmd := exec.Command(args[0], args[1:]...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v failed: %s\n%s", args, err, out)
		}
	}
	os.WriteFile(filepath.Join(dir, ".gitkeep"), []byte(""), 0644)
	for _, args := range [][]string{
		{"git", "-C", dir, "add", ".gitkeep"},
		{"git", "-C", dir, "commit", "-m", "initial"},
	} {
		cmd := exec.Command(args[0], args[1:]...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v failed: %s\n%s", args, err, out)
		}
	}
	return dir
}

func newTestStore(t *testing.T) *internal.RepositoryBranches {
	t.Helper()
	dir := t.TempDir()
	return &internal.RepositoryBranches{
		RepositoryName: "test-repo",
		StoreDirectory: dir,
	}
}

func chdir(t *testing.T, dir string) {
	t.Helper()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	t.Cleanup(func() { os.Chdir(origDir) })
}

func TestExecuteCreate_NormalBranch(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	executeCreate(&git, store, "feature/test", "ft", "test note ", createOpts{})

	// Verify branch exists in git and is checked out
	g2 := internal.Git{}
	current, _ := g2.GetCurrentBranch()
	if current != "feature/test" {
		t.Errorf("expected checkout to 'feature/test', got %q", current)
	}

	// Verify registered in store
	if !store.BranchExists(internal.Branch{Name: "feature/test"}) {
		t.Error("branch should be registered")
	}
	b, _ := store.GetBranchByAlias("ft")
	if b.Note != "test note " {
		t.Errorf("expected note 'test note ', got %q", b.Note)
	}
}

func TestExecuteCreate_OnlyCreate(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	beforeBranch, _ := git.GetCurrentBranch()

	executeCreate(&git, store, "feature/no-checkout", "nc", "", createOpts{CreateOnly: true})

	// Should still be on original branch
	g2 := internal.Git{}
	afterBranch, _ := g2.GetCurrentBranch()
	if afterBranch != beforeBranch {
		t.Errorf("should stay on %q with CreateOnly, got %q", beforeBranch, afterBranch)
	}

	// But branch should be registered
	if !store.BranchExists(internal.Branch{Name: "feature/no-checkout"}) {
		t.Error("branch should be registered even with CreateOnly")
	}
}

func TestExecuteCreate_WithPrefix(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	executeCreate(&git, store, "my-feature", "mf", "", createOpts{Prefix: "dev/"})

	// Branch name should have prefix applied
	if !store.BranchExists(internal.Branch{Name: "dev/my-feature"}) {
		t.Error("expected branch name with prefix 'dev/my-feature'")
	}
	b, found := store.GetBranchByAlias("mf")
	if !found || b.Name != "dev/my-feature" {
		t.Errorf("expected prefixed name, got %+v", b)
	}
}

func TestExecuteCreate_WithWorktree(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	wtDir := filepath.Join(repoDir, "worktrees", "wt1")

	executeCreate(&git, store, "feature/wt", "wt", "note ", createOpts{
		UseWorktree:          true,
		WorktreeAlias:        "wt1",
		WorktreePathTemplate: "./worktrees/{alias}",
		RepoPath:             repoDir,
		RepoName:             filepath.Base(repoDir),
	})

	// Verify worktree directory exists
	if _, err := os.Stat(wtDir); os.IsNotExist(err) {
		t.Error("worktree directory should exist")
	}

	// Verify both registered
	if !store.BranchExists(internal.Branch{Name: "feature/wt"}) {
		t.Error("branch should be registered")
	}
	if !store.WorktreeAliasExists("wt1") {
		t.Error("worktree should be registered")
	}
}

func TestExecuteCreate_WorktreeDefaultAlias(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	store := newTestStore(t)
	git := internal.Git{}

	// When WorktreeAlias is empty, it should default to the branch alias
	executeCreate(&git, store, "feature/def", "defal", "", createOpts{
		UseWorktree:          true,
		WorktreeAlias:        "",
		WorktreePathTemplate: "./worktrees/{alias}",
		RepoPath:             repoDir,
		RepoName:             filepath.Base(repoDir),
	})

	// Worktree alias should default to branch alias "defal"
	if !store.WorktreeAliasExists("defal") {
		t.Error("worktree alias should default to branch alias")
	}
}
