package cmd

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/spf13/cobra"
)

// promptStringCmd represents the promptString command
var promptStringCmd = &cobra.Command{
	Use:    "_ps",
	Short:  "Returns the prompt string - used for the custom prompt",
	Long:   `Returns the prompt string - used for the custom prompt`,
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		prompt := ""
		if runtime.GOOS == "windows" {
			username := os.Getenv("USERNAME")
			prompt = fmt.Sprintf("%s / %s # ", username, getPrompt(" > "))
		} else {
			username := os.Getenv("LOGNAME")
			if username == "" {
				currentUser, err := user.Current()
				if err == nil {
					username = currentUser.Username
				}
			}
			prompt = fmt.Sprintf("%s ➤ %s ❖ ", username, getPrompt(" ⌥ "))
		}
		fmt.Print(prompt)
	},
}

// Build the custom prompt string
func getPrompt(separator string) string {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return "$"
	}

	git := internal.Git{}
	repo, err := git.GetRepositoryName()
	if err != nil || repo == "" {
		return workingDirectory
	}

	currentBranch, err := git.GetCurrentBranch()
	if err != nil || currentBranch == "" {
		return workingDirectory
	}

	branchAlias := ""
	repoBranches := internal.GetRepositoryBranches()
	branch, ok := repoBranches.GetBranchByName(currentBranch)
	if ok {
		branchAlias = branch.Alias
	}

	return buildPrompt(repo, currentBranch, branchAlias, workingDirectory, separator)
}

// buildPrompt assembles the prompt string from resolved values.
func buildPrompt(repoName, currentBranch, branchAlias, workingDir, separator string) string {
	// Format branch with alias if present
	branchDisplay := currentBranch
	if branchAlias != "" {
		branchDisplay = fmt.Sprintf("%s (%s)", currentBranch, branchAlias)
	}

	// Compute subpath relative to repo
	subpath := ""
	index := strings.Index(workingDir, repoName)
	if index >= 0 {
		subpath = workingDir[index+len(repoName):]
		subpath = strings.TrimPrefix(subpath, "/")
		if subpath != "" {
			subpath = fmt.Sprintf(" [%s]", subpath)
		}
	}

	return repoName + subpath + separator + branchDisplay
}

func init() {
	rootCmd.AddCommand(promptStringCmd)
}
