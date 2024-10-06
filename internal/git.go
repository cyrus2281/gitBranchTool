package internal

import (
	"fmt"
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
	output, err := cmd.CombinedOutput()
	if err != nil {
		if len(output) > 0 {
			return "", fmt.Errorf("%v", strings.TrimSpace(string(output)))
		}
		return "", fmt.Errorf("%v", output)
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

func (g *Git) SwitchBranch(name string) error {
	_, err := runCommand([]string{"git", "checkout", name})
	return err
}

func (g *Git) GetBranches() ([]string, error) {
	output, err := runCommand([]string{"git", "branch"})
	if err != nil {
		return nil, err
	}
	branches := strings.Split(output, "\n")
	parsedBranches := make([]string, 0, len(branches))
	for _, branch := range branches {
		branch = strings.TrimSpace(branch)
		if branch == "" {
			continue
		}
		if strings.HasPrefix(branch, "*") {
			branch = strings.TrimSpace(branch[1:])
		}
		parsedBranches = append(parsedBranches, branch)
	}
	return parsedBranches, nil
}
