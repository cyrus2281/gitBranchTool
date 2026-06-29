package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cyrus2281/gitBranchTool/internal"
)

// runGitIn runs a git command in the given dir and fails the test on error.
func runGitIn(t *testing.T, dir string, args ...string) string {
	t.Helper()
	full := append([]string{"-C", dir}, args...)
	out, err := exec.Command("git", full...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %s\n%s", args, err, out)
	}
	return string(out)
}

// commitFileIn writes content to file in dir and commits it.
func commitFileIn(t *testing.T, dir, file, content, msg string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, file), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	runGitIn(t, dir, "add", file)
	runGitIn(t, dir, "commit", "-m", msg)
}

func currentBranch(t *testing.T, dir string) string {
	t.Helper()
	return strings.TrimSpace(runGitIn(t, dir, "branch", "--show-current"))
}

func TestValidateMergeOpts(t *testing.T) {
	cases := []struct {
		name    string
		opts    mergeOpts
		hasArg  bool
		wantErr bool
	}{
		{"valid merge", mergeOpts{}, true, false},
		{"valid rebase", mergeOpts{Rebase: true}, true, false},
		{"merge squash+noverify", mergeOpts{Squash: true, NoVerify: true}, true, false},
		{"abort no arg", mergeOpts{Abort: true}, false, false},
		{"continue no arg", mergeOpts{Continue: true}, false, false},
		{"abort+continue", mergeOpts{Abort: true, Continue: true}, false, true},
		{"abort with arg", mergeOpts{Abort: true}, true, true},
		{"continue with arg", mergeOpts{Continue: true}, true, true},
		{"rebase+squash", mergeOpts{Rebase: true, Squash: true}, true, true},
		{"rebase+noverify", mergeOpts{Rebase: true, NoVerify: true}, true, true},
		{"rebase+ffonly", mergeOpts{Rebase: true, FFOnly: true}, true, true},
		{"rebase+noff", mergeOpts{Rebase: true, NoFF: true}, true, true},
		{"ffonly+noff", mergeOpts{FFOnly: true, NoFF: true}, true, true},
	}
	for _, c := range cases {
		err := validateMergeOpts(c.opts, c.hasArg)
		if (err != nil) != c.wantErr {
			t.Errorf("%s: wantErr=%v, got %v", c.name, c.wantErr, err)
		}
	}
}

func TestExecuteMerge_MergeByAlias(t *testing.T) {
	repoDir := initTestRepo(t)
	def := currentBranch(t, repoDir)
	runGitIn(t, repoDir, "checkout", "-b", "feature/x")
	commitFileIn(t, repoDir, "feat.txt", "hi", "feat commit")
	runGitIn(t, repoDir, "checkout", def)
	chdir(t, repoDir)
	t.Setenv("GIT_EDITOR", "true")

	store := newTestStore(t)
	store.AddBranch(internal.Branch{Name: "feature/x", Alias: "fx"})
	git := internal.Git{}

	executeMerge(&git, store, "fx", mergeOpts{})

	if _, err := os.Stat(filepath.Join(repoDir, "feat.txt")); err != nil {
		t.Errorf("expected feat.txt merged into current branch: %v", err)
	}
}

func TestExecuteMerge_Rebase(t *testing.T) {
	repoDir := initTestRepo(t)
	def := currentBranch(t, repoDir)
	runGitIn(t, repoDir, "checkout", "-b", "feature")
	commitFileIn(t, repoDir, "f.txt", "f", "f commit")
	runGitIn(t, repoDir, "checkout", def)
	commitFileIn(t, repoDir, "d.txt", "d", "d commit")
	runGitIn(t, repoDir, "checkout", "feature")
	chdir(t, repoDir)

	store := newTestStore(t)
	git := internal.Git{}

	executeMerge(&git, store, def, mergeOpts{Rebase: true})

	for _, f := range []string{"d.txt", "f.txt"} {
		if _, err := os.Stat(filepath.Join(repoDir, f)); err != nil {
			t.Errorf("expected %s present after rebase: %v", f, err)
		}
	}
}

func TestExecuteMerge_Squash(t *testing.T) {
	repoDir := initTestRepo(t)
	def := currentBranch(t, repoDir)
	runGitIn(t, repoDir, "checkout", "-b", "topic")
	commitFileIn(t, repoDir, "s.txt", "s", "s commit")
	runGitIn(t, repoDir, "checkout", def)
	chdir(t, repoDir)

	store := newTestStore(t)
	git := internal.Git{}

	executeMerge(&git, store, "topic", mergeOpts{Squash: true})

	if _, err := os.Stat(filepath.Join(repoDir, "s.txt")); err != nil {
		t.Errorf("expected s.txt in working tree after squash: %v", err)
	}
	if staged := runGitIn(t, repoDir, "diff", "--cached", "--name-only"); !strings.Contains(staged, "s.txt") {
		t.Errorf("expected s.txt to be staged, got %q", staged)
	}
	if log := runGitIn(t, repoDir, "log", "--oneline"); strings.Contains(log, "s commit") {
		t.Error("squash should not bring the topic commit into history before commit")
	}
}

func TestExecuteMerge_DefaultBranch(t *testing.T) {
	repoDir := initTestRepo(t)
	def := currentBranch(t, repoDir)
	// Branch "work" off the initial commit, then advance the default branch.
	runGitIn(t, repoDir, "checkout", "-b", "work")
	runGitIn(t, repoDir, "checkout", def)
	commitFileIn(t, repoDir, "main-new.txt", "m", "main advance")
	runGitIn(t, repoDir, "checkout", "work")
	chdir(t, repoDir)
	t.Setenv("GIT_EDITOR", "true")

	store := newTestStore(t)
	store.SetDefaultBranch(def)
	git := internal.Git{}

	// Simulates the no-arg path in Run: id = store.GetDefaultBranch().
	executeMerge(&git, store, store.GetDefaultBranch(), mergeOpts{})

	if _, err := os.Stat(filepath.Join(repoDir, "main-new.txt")); err != nil {
		t.Errorf("expected default branch merged into work: %v", err)
	}
}

func TestExecuteMerge_SameBranchGuard(t *testing.T) {
	repoDir := initTestRepo(t)
	def := currentBranch(t, repoDir)
	chdir(t, repoDir)

	store := newTestStore(t)
	git := internal.Git{}

	before := strings.TrimSpace(runGitIn(t, repoDir, "rev-parse", "HEAD"))
	executeMerge(&git, store, def, mergeOpts{})
	after := strings.TrimSpace(runGitIn(t, repoDir, "rev-parse", "HEAD"))

	if before != after {
		t.Errorf("same-branch guard should not change HEAD: %q -> %q", before, after)
	}
}

func TestExecuteMerge_Fetch(t *testing.T) {
	repoDir := initTestRepo(t)
	def := currentBranch(t, repoDir)

	remoteDir := t.TempDir()
	if out, err := exec.Command("git", "init", "--bare", remoteDir).CombinedOutput(); err != nil {
		t.Fatalf("git init --bare failed: %s\n%s", err, out)
	}
	runGitIn(t, repoDir, "remote", "add", "origin", remoteDir)

	// Push topic@topic1, then push topic2, then rewind the local branch so the
	// remote is strictly ahead of the local tip.
	runGitIn(t, repoDir, "checkout", "-b", "topic")
	commitFileIn(t, repoDir, "topic1.txt", "1", "topic1")
	runGitIn(t, repoDir, "push", "origin", "topic")
	commitFileIn(t, repoDir, "topic2.txt", "2", "topic2")
	runGitIn(t, repoDir, "push", "origin", "topic")
	runGitIn(t, repoDir, "reset", "--hard", "HEAD~1") // local topic now at topic1
	runGitIn(t, repoDir, "checkout", def)
	chdir(t, repoDir)
	t.Setenv("GIT_EDITOR", "true")

	store := newTestStore(t)
	git := internal.Git{}

	executeMerge(&git, store, "topic", mergeOpts{Fetch: true, Remote: "origin"})

	// topic2.txt only exists on the remote tip; its presence proves we merged
	// the freshly fetched remote branch rather than the stale local one.
	if _, err := os.Stat(filepath.Join(repoDir, "topic2.txt")); err != nil {
		t.Errorf("expected topic2.txt from fetched remote tip: %v", err)
	}
}

// setupConflict creates a merge conflict on c.txt between "feature" and the
// default branch, leaving the repo checked out on the default branch.
func setupConflict(t *testing.T, repoDir string) string {
	t.Helper()
	def := currentBranch(t, repoDir)
	commitFileIn(t, repoDir, "c.txt", "base\n", "base")
	runGitIn(t, repoDir, "checkout", "-b", "feature")
	commitFileIn(t, repoDir, "c.txt", "feature\n", "feature change")
	runGitIn(t, repoDir, "checkout", def)
	commitFileIn(t, repoDir, "c.txt", "default\n", "default change")
	return def
}

// TestMergeSubprocessHelper is invoked as a subprocess by the exit-path tests
// below so we can observe os.Exit from a conflict / fatal.
func TestMergeSubprocessHelper(t *testing.T) {
	if os.Getenv("MERGE_TEST_SUBPROCESS") != "1" {
		t.Skip("subprocess helper — invoked by exit-path tests")
	}
	repoDir := os.Getenv("MERGE_TEST_REPO_DIR")
	os.Chdir(repoDir)

	store := &internal.RepositoryBranches{
		RepositoryName: filepath.Base(repoDir),
		StoreDirectory: os.Getenv("MERGE_TEST_STORE_DIR"),
	}
	git := &internal.Git{}

	opts := mergeOpts{Remote: "origin"}
	opts.Rebase = os.Getenv("MERGE_TEST_REBASE") == "1"
	opts.Abort = os.Getenv("MERGE_TEST_ABORT") == "1"
	opts.Continue = os.Getenv("MERGE_TEST_CONTINUE") == "1"

	if os.Getenv("MERGE_TEST_MODE") == "state" {
		executeMergeState(git, opts)
		return
	}
	executeMerge(git, store, os.Getenv("MERGE_TEST_ID"), opts)
}

func TestExecuteMerge_ConflictExitsAndAbort(t *testing.T) {
	repoDir := initTestRepo(t)
	setupConflict(t, repoDir)
	storeDir := t.TempDir()

	cmd := exec.Command(os.Args[0], "-test.run=^TestMergeSubprocessHelper$")
	cmd.Dir = repoDir
	cmd.Env = append(os.Environ(),
		"MERGE_TEST_SUBPROCESS=1",
		"MERGE_TEST_MODE=merge",
		"MERGE_TEST_REPO_DIR="+repoDir,
		"MERGE_TEST_STORE_DIR="+storeDir,
		"MERGE_TEST_ID=feature",
		"GIT_EDITOR=true",
	)
	if err := cmd.Run(); err == nil {
		t.Fatal("expected merge conflict to exit non-zero")
	}

	chdir(t, repoDir)
	git := internal.Git{}
	if !git.MergeInProgress() {
		t.Fatal("expected merge in progress after conflict")
	}
	executeMergeState(&git, mergeOpts{Abort: true})
	if git.MergeInProgress() {
		t.Error("expected merge to be aborted")
	}
}

func TestExecuteMergeState_Continue(t *testing.T) {
	repoDir := initTestRepo(t)
	setupConflict(t, repoDir)
	chdir(t, repoDir)
	t.Setenv("GIT_EDITOR", "true")

	git := internal.Git{}
	if err := git.MergeBranch("feature", nil); err == nil {
		t.Fatal("expected merge to conflict")
	}
	if !git.MergeInProgress() {
		t.Fatal("expected merge in progress")
	}

	// Resolve the conflict, then continue.
	if err := os.WriteFile(filepath.Join(repoDir, "c.txt"), []byte("resolved\n"), 0644); err != nil {
		t.Fatal(err)
	}
	runGitIn(t, repoDir, "add", "c.txt")

	executeMergeState(&git, mergeOpts{Continue: true})

	if git.MergeInProgress() {
		t.Error("expected merge to complete after --continue")
	}
}

func TestExecuteMergeState_NothingInProgress(t *testing.T) {
	repoDir := initTestRepo(t)
	storeDir := t.TempDir()

	cmd := exec.Command(os.Args[0], "-test.run=^TestMergeSubprocessHelper$")
	cmd.Dir = repoDir
	cmd.Env = append(os.Environ(),
		"MERGE_TEST_SUBPROCESS=1",
		"MERGE_TEST_MODE=state",
		"MERGE_TEST_REPO_DIR="+repoDir,
		"MERGE_TEST_STORE_DIR="+storeDir,
		"MERGE_TEST_CONTINUE=1",
	)
	if err := cmd.Run(); err == nil {
		t.Fatal("expected --continue with nothing in progress to exit non-zero")
	}
}
