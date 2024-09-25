package internal

import (
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

var gitCommands = map[string][]string{
	"currentBranch":  {"git", "branch", "--show-current"},
	"repositoryPath": {"git", "rev-parse", "--show-toplevel"},
	"isGitRepo":      {"git", "rev-parse", "--is-inside-work-tree"},
}

func runCommand(command []string) (string, error) {
	// Run the command
	cmd := exec.Command(command[0], command[1:]...)

	// Run the command and capture the output
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to run git command: %v", err)
		return "", err
	}

	// Convert the output to a string and trim any whitespace
	return strings.TrimSpace(string(output)), nil
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

// GetRepositoryName returns the name of the git repository
func (g *Git) GetRepositoryName() (string, error) {
	if g.repositoryName != "" {
		return g.repositoryName, nil
	}
	// Run the command
	repositoryPath, err := runCommand(gitCommands["repositoryPath"])
	if err != nil {
		return "", err
	}
	// get the basename of the path
	g.repositoryName = filepath.Base(repositoryPath)

	return g.repositoryName, nil
}
