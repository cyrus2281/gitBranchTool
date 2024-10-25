/*
Copyright © 2024 Cyrus Mobini
*/
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
	// Repository name
	repo, err := git.GetRepositoryName()
	if err != nil || repo == "" {
		return workingDirectory
	}

	// Branch Name
	currentBranch, err := git.GetCurrentBranch()
	if err != nil || currentBranch == "" {
		return workingDirectory
	}

	// Alias Name
	repoBranches := internal.GetRepositoryBranches()
	branch, ok := repoBranches.GetBranchByName(currentBranch)
	if ok && branch.Alias != "" {
		currentBranch = fmt.Sprintf("%s (%s)", currentBranch, branch.Alias)
	}

	// Subpath
	subpath := ""
	index := strings.Index(workingDirectory, repo)
	if index >= 0 {
		subpath = workingDirectory[index+len(repo):]
		subpath = strings.TrimPrefix(subpath, "/")
		if subpath != "" {
			subpath = fmt.Sprintf(" [%s]", subpath)
		}
	}

	// Custom Prompt String
	return repo + subpath + separator + currentBranch
}

func init() {
	rootCmd.AddCommand(promptStringCmd)
}
