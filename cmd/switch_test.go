package cmd

import (
	"os"
	"os/exec"
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

// TestSwitchSubprocessHelper is invoked as a subprocess by the fail tests below.
// It calls executeSwitch directly so we can observe os.Exit from logger.Fatalln.
func TestSwitchSubprocessHelper(t *testing.T) {
	if os.Getenv("SWITCH_TEST_SUBPROCESS") != "1" {
		t.Skip("subprocess helper — invoked by fail tests")
	}

	repoDir := os.Getenv("SWITCH_TEST_REPO_DIR")
	storeDir := os.Getenv("SWITCH_TEST_STORE_DIR")
	os.Chdir(repoDir)

	store := &internal.RepositoryBranches{
		RepositoryName: filepath.Base(repoDir),
		StoreDirectory: storeDir,
	}
	git := &internal.Git{}

	var opts switchOpts
	if os.Getenv("SWITCH_TEST_WORKTREE") == "1" {
		repoPath, _ := git.GetRepositoryPath()
		opts = switchOpts{
			UseWorktree:          true,
			WorktreePathTemplate: "./worktrees/{alias}",
			RepoPath:             repoPath,
			RepoName:             filepath.Base(repoDir),
		}
	}

	executeSwitch(git, store,
		os.Getenv("SWITCH_TEST_ID"),
		os.Getenv("SWITCH_TEST_ALIAS"),
		"",
		opts,
	)
}

// TestExecuteSwitch_FailNoRegister verifies that a failed branch switch
// (nonexistent branch) does not register the branch in the store.
func TestExecuteSwitch_FailNoRegister(t *testing.T) {
	repoDir := initTestRepo(t)
	storeDir := t.TempDir()
	storeDir, _ = filepath.EvalSymlinks(storeDir)

	cmd := exec.Command(os.Args[0], "-test.run=^TestSwitchSubprocessHelper$")
	cmd.Dir = repoDir
	cmd.Env = append(os.Environ(),
		"SWITCH_TEST_SUBPROCESS=1",
		"SWITCH_TEST_REPO_DIR="+repoDir,
		"SWITCH_TEST_STORE_DIR="+storeDir,
		"SWITCH_TEST_ID=nonexistent-branch",
		"SWITCH_TEST_ALIAS=nb",
	)
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected switch to nonexistent branch to fail")
	}

	store := &internal.RepositoryBranches{
		RepositoryName: filepath.Base(repoDir),
		StoreDirectory: storeDir,
	}
	if len(store.GetBranches()) != 0 {
		t.Errorf("expected 0 branches after failed switch, got %d", len(store.GetBranches()))
	}
}

// TestExecuteSwitch_WorktreeMode_FailNoRegister verifies that a failed worktree
// creation (nonexistent branch) registers neither branch nor worktree.
func TestExecuteSwitch_WorktreeMode_FailNoRegister(t *testing.T) {
	repoDir := initTestRepo(t)
	storeDir := t.TempDir()
	storeDir, _ = filepath.EvalSymlinks(storeDir)

	cmd := exec.Command(os.Args[0], "-test.run=^TestSwitchSubprocessHelper$")
	cmd.Dir = repoDir
	cmd.Env = append(os.Environ(),
		"SWITCH_TEST_SUBPROCESS=1",
		"SWITCH_TEST_REPO_DIR="+repoDir,
		"SWITCH_TEST_STORE_DIR="+storeDir,
		"SWITCH_TEST_ID=nonexistent-branch",
		"SWITCH_TEST_ALIAS=nb",
		"SWITCH_TEST_WORKTREE=1",
	)
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected worktree creation for nonexistent branch to fail")
	}

	store := &internal.RepositoryBranches{
		RepositoryName: filepath.Base(repoDir),
		StoreDirectory: storeDir,
	}
	if len(store.GetBranches()) != 0 {
		t.Errorf("expected 0 branches after failed worktree creation, got %d", len(store.GetBranches()))
	}
	if len(store.GetWorktrees()) != 0 {
		t.Errorf("expected 0 worktrees after failed worktree creation, got %d", len(store.GetWorktrees()))
	}
}
