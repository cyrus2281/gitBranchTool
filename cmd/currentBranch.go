/*
Copyright © 2024 Cyrus Mobini
*/
package cmd

import (
	"fmt"

	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/cyrus2281/go-logger"
	"github.com/spf13/cobra"
)

// currentBranchCmd represents the currentBranch command
var currentBranchCmd = &cobra.Command{
	Use:     "currentBranch",
	Short:   "Returns the name of active branch with alias and note",
	Long:    `Returns the name of active branch with alias and note`,
	Aliases: []string{"current-branch", "cb"},
	Run: func(cmd *cobra.Command, args []string) {
		git := internal.Git{}
		if !git.IsGitRepo() {
			logger.Fatalln("Not a git repository")
		}

		currentBranch, err := git.GetCurrentBranch()
		logger.CheckFatalln(err)

		repoBranches := internal.GetRepositoryBranches()
		branch, ok := repoBranches.GetBranchByName(currentBranch)
		if !ok {
			logger.FatalF("Branch \"%v\" is not registered\n", currentBranch)
		}
		internal.PrintTableHeader()
		fmt.Println(branch.String())
	},
}

func init() {
	rootCmd.AddCommand(currentBranchCmd)
}
