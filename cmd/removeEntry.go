/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/cyrus2281/go-logger"
	"github.com/spf13/cobra"
)

// removeEntryCmd represents the removeEntry command
var removeEntryCmd = &cobra.Command{
	Use:     "removeEntry NAME|ALIAS",
	Short:   "Removes a registered branch entry without deleting the branch",
	Long:    `Removes a registered branch entry without deleting the branch",`,
	Aliases: []string{"remove-entry", "re"},
	Run: func(cmd *cobra.Command, args []string) {
		git := internal.Git{}
		if !git.IsGitRepo() {
			logger.Fatalln("Not a git repository")
		}

		id := args[0]
		repoBranches := internal.GetRepositoryBranches()
		branch, ok := repoBranches.GetBranchByNameOrAlias(id)
		if !ok {
			logger.Fatalln("Branch/Alias is not registered")
		}
		repoBranches.RemoveBranch(branch)
		logger.InfoF("Branch \"%v\" with alias \"%v\" was removed\n", branch.Name, branch.Alias)
	},
}

func init() {
	rootCmd.AddCommand(removeEntryCmd)
}
