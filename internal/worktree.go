package internal

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Worktree struct {
	Alias string `json:"alias"`
	Path  string `json:"path"`
	Note  string `json:"note"`
}

// Convert to print format: Path | Alias | Branch | Branch Alias | Note
// Branch and Branch Alias are populated at display time, not stored
func (w *Worktree) String() string {
	return fmt.Sprintf("%-40s\t%-15s\t%-25s\t%-15s\t%-20s", w.Path, w.Alias, "", "", w.Note)
}

// StringWithBranch returns formatted string with branch info populated
func (w *Worktree) StringWithBranch(branch string, branchAlias string) string {
	return fmt.Sprintf("%-40s\t%-15s\t%-25s\t%-15s\t%-20s", w.Path, w.Alias, branch, branchAlias, w.Note)
}

// ResolveWorktreePath replaces template variables and resolves to absolute path
// Template variables: {repository}, {alias}, {branch}
func ResolveWorktreePath(template, repoRoot, repoName, alias, branch string) string {
	result := template
	result = strings.ReplaceAll(result, "{repository}", repoName)
	result = strings.ReplaceAll(result, "{alias}", alias)
	result = strings.ReplaceAll(result, "{branch}", branch)

	// If the path is relative, resolve it relative to the repo root
	if !filepath.IsAbs(result) {
		result = filepath.Join(repoRoot, result)
	}

	// Clean the path
	result = filepath.Clean(result)
	return result
}

// ParseWorktreeList parses the output of `git worktree list --porcelain`
// Returns a map of worktree path to branch name
func ParseWorktreeList(output string) map[string]string {
	result := make(map[string]string)
	if output == "" {
		return result
	}

	lines := strings.Split(output, "\n")
	currentPath := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "worktree ") {
			currentPath = strings.TrimPrefix(line, "worktree ")
			if currentPath != "" {
				// Record immediately so detached worktrees (no branch line) are preserved
				result[currentPath] = ""
			}
		} else if strings.HasPrefix(line, "branch ") {
			branchRef := strings.TrimPrefix(line, "branch ")
			// Convert refs/heads/branch-name to branch-name
			branchName := strings.TrimPrefix(branchRef, "refs/heads/")
			if currentPath != "" {
				result[currentPath] = branchName
			}
		} else if line == "" {
			currentPath = ""
		}
	}

	return result
}

// GetWorktreePathForBranch returns the worktree path for a given branch name
// Returns empty string if no worktree is found for the branch
func GetWorktreePathForBranch(worktreeMap map[string]string, branchName string) string {
	for path, branch := range worktreeMap {
		if branch == branchName {
			return path
		}
	}
	return ""
}
