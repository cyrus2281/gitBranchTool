package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cyrus2281/go-logger"
)

var gitCommands = map[string][]string{
	"currentBranch":  {"git", "branch", "--show-current"},
	"repositoryPath": {"git", "rev-parse", "--show-toplevel"},
	"gitCommonDir":   {"git", "rev-parse", "--git-common-dir"},
	"isGitRepo":      {"git", "rev-parse", "--is-inside-work-tree"},
}

func runCommand(command []string) (string, error) {
	logger.Debugln("Running:", strings.Join(command, " "))
	// Run the command
	cmd := exec.Command(command[0], command[1:]...)

	// Run the command and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		if len(output) > 0 {
			return "", fmt.Errorf("%v", strings.TrimSpace(string(output)))
		}
		return "", err
	}

	// Convert the output to a string and trim any whitespace
	return strings.TrimSpace(string(output)), nil
}

// runInteractiveCommand attaches the user's terminal so interactive git
// operations (commit-message editor, live conflict output) behave like
// native git, instead of capturing and swallowing their output.
func runInteractiveCommand(command []string) error {
	logger.Debugln("Running:", strings.Join(command, " "))
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type Git struct {
	currentBranch  string
	repositoryName string
}

// IsGitRepo returns true if the current directory is inside a git repository
func (g *Git) IsGitRepo() bool {
	_, err := runCommand(gitCommands["isGitRepo"])
	return err == nil
}

// GetCurrentBranch returns the current branch of the git repository
func (g *Git) GetCurrentBranch() (string, error) {
	if g.currentBranch != "" {
		return g.currentBranch, nil
	}
	currentBranch, err := runCommand(gitCommands["currentBranch"])
	if err != nil {
		return "", err
	}
	g.currentBranch = currentBranch
	return g.currentBranch, nil
}

// GetRepositoryName returns the name of the main git repository.
// Uses --git-common-dir to resolve the correct repo even from a worktree.
func (g *Git) GetRepositoryName() (string, error) {
	if g.repositoryName != "" {
		return g.repositoryName, nil
	}
	// Use --git-common-dir to get the main repo's .git path,
	// which works correctly even when inside a worktree.
	gitCommonDir, err := runCommand(gitCommands["gitCommonDir"])
	if err != nil {
		return "", err
	}
	// gitCommonDir returns something like:
	//   ".git" (when in the main repo)
	//   "/absolute/path/to/main-repo/.git" (when in a worktree)
	//   "/absolute/path/to/main-repo/.git" (bare worktrees)
	// We need the parent of the .git directory to get the repo name.
	if !filepath.IsAbs(gitCommonDir) {
		// Relative path (e.g. ".git") — resolve relative to --show-toplevel
		repositoryPath, err := runCommand(gitCommands["repositoryPath"])
		if err != nil {
			return "", err
		}
		g.repositoryName = filepath.Base(repositoryPath)
	} else {
		// Absolute path — parent of .git dir is the main repo
		g.repositoryName = filepath.Base(filepath.Dir(gitCommonDir))
	}

	return g.repositoryName, nil
}

// CreateNewBranch creates a new branch
func (g *Git) CreateNewBranch(name string) error {
	_, err := runCommand([]string{"git", "branch", name})
	return err
}

// SwitchToNewBranch creates a new branch and switches to it
func (g *Git) SwitchToNewBranch(name string) error {
	_, err := runCommand([]string{"git", "checkout", "-b", name})
	return err
}

func (g *Git) DeleteBranch(name string, force bool) error {
	if force {
		_, err := runCommand([]string{"git", "branch", "-D", name})
		return err
	} else {
		_, err := runCommand([]string{"git", "branch", "-d", name})
		return err
	}
}

func (g *Git) DeleteRemoteBranch(name string) error {
	_, err := runCommand([]string{"git", "push", "origin", "--delete", name})
	return err
}

func (g *Git) SwitchBranch(name string) error {
	_, err := runCommand([]string{"git", "checkout", name})
	return err
}

// MergeBranch merges ref into the current branch. opts are extra git merge
// flags (e.g. --squash, --no-verify, --no-ff). Runs interactively so conflict
// output and the merge-commit editor work as with native git.
func (g *Git) MergeBranch(ref string, opts []string) error {
	command := append([]string{"git", "merge"}, opts...)
	command = append(command, ref)
	return runInteractiveCommand(command)
}

// RebaseBranch rebases the current branch onto ref.
func (g *Git) RebaseBranch(ref string) error {
	return runInteractiveCommand([]string{"git", "rebase", ref})
}

// MergeAbort aborts an in-progress merge.
func (g *Git) MergeAbort() error {
	_, err := runCommand([]string{"git", "merge", "--abort"})
	return err
}

// MergeContinue continues an in-progress merge after conflicts are resolved.
func (g *Git) MergeContinue() error {
	return runInteractiveCommand([]string{"git", "merge", "--continue"})
}

// RebaseAbort aborts an in-progress rebase.
func (g *Git) RebaseAbort() error {
	_, err := runCommand([]string{"git", "rebase", "--abort"})
	return err
}

// RebaseContinue continues an in-progress rebase after conflicts are resolved.
func (g *Git) RebaseContinue() error {
	return runInteractiveCommand([]string{"git", "rebase", "--continue"})
}

// Fetch fetches the given branch from the given remote. The fetched tip is
// available afterwards as FETCH_HEAD.
func (g *Git) Fetch(remote, branch string) error {
	return runInteractiveCommand([]string{"git", "fetch", remote, branch})
}

// GetGitDir returns the path to the git directory for the current worktree.
// Inside a worktree this is the worktree-specific git dir, which is where
// merge/rebase state files live.
func (g *Git) GetGitDir() (string, error) {
	return runCommand([]string{"git", "rev-parse", "--git-dir"})
}

// MergeInProgress reports whether a merge is currently in progress.
func (g *Git) MergeInProgress() bool {
	gitDir, err := g.GetGitDir()
	if err != nil {
		return false
	}
	_, err = os.Stat(filepath.Join(gitDir, "MERGE_HEAD"))
	return err == nil
}

// RebaseInProgress reports whether a rebase is currently in progress.
func (g *Git) RebaseInProgress() bool {
	gitDir, err := g.GetGitDir()
	if err != nil {
		return false
	}
	for _, name := range []string{"rebase-merge", "rebase-apply"} {
		if _, err := os.Stat(filepath.Join(gitDir, name)); err == nil {
			return true
		}
	}
	return false
}

// WorktreeAdd creates a worktree for an existing branch
func (g *Git) WorktreeAdd(path, branch string) error {
	_, err := runCommand([]string{"git", "worktree", "add", path, branch})
	return err
}

// WorktreeAddNewBranch creates a new branch and a worktree for it
func (g *Git) WorktreeAddNewBranch(branch, path string) error {
	_, err := runCommand([]string{"git", "worktree", "add", "-b", branch, path})
	return err
}

// WorktreeRemove removes a worktree by path
func (g *Git) WorktreeRemove(path string, force bool) error {
	if force {
		_, err := runCommand([]string{"git", "worktree", "remove", path, "--force"})
		return err
	}
	_, err := runCommand([]string{"git", "worktree", "remove", path})
	return err
}

// WorktreeList returns the porcelain output of git worktree list
func (g *Git) WorktreeList() (string, error) {
	return runCommand([]string{"git", "worktree", "list", "--porcelain"})
}

// WorktreePrune removes stale worktree entries
func (g *Git) WorktreePrune() error {
	_, err := runCommand([]string{"git", "worktree", "prune"})
	return err
}

// GetRepositoryPath returns the absolute path to the repository root
func (g *Git) GetRepositoryPath() (string, error) {
	return runCommand(gitCommands["repositoryPath"])
}

func (g *Git) GetBranches() ([]string, error) {
	output, err := runCommand([]string{"git", "branch"})
	if err != nil {
		return nil, err
	}
	return parseBranchOutput(output), nil
}

// parseBranchOutput parses the output of `git branch` into a list of branch names.
// It strips the leading marker git prepends to a branch (`*` for the current
// branch, `+` for a branch checked out in a linked worktree), trims whitespace,
// and filters out detached HEAD entries.
func parseBranchOutput(output string) []string {
	lines := strings.Split(output, "\n")
	parsedBranches := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "*") || strings.HasPrefix(line, "+") {
			line = strings.TrimSpace(line[1:])
		}
		// Filter out detached HEAD entries like "(HEAD detached at abc123)"
		// Uses Contains rather than HasPrefix("(") to avoid dropping valid
		// branch names that happen to start with a parenthesis.
		if strings.Contains(line, "HEAD detached") {
			continue
		}
		parsedBranches = append(parsedBranches, line)
	}
	return parsedBranches
}
