/*
Copyright © 2024 Cyrus Mobini
*/
package cmd

import (
	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/cyrus2281/go-logger"
	"github.com/spf13/cobra"
)

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:     "rename NAME ALIAS",
	Short:   "Updates the alias for the given branch name",
	Long:    `Updates the alias for the given branch name.`,
	Args:    cobra.ExactArgs(2),
	Aliases: []string{"mv", "updateBranchAlias", "update-branch-alias", "uba"},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		position := len(args) + 1
		if position == 1 {
			return internal.GetBranches()
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		git := internal.Git{}
		if !git.IsGitRepo() {
			logger.Fatalln("Not a git repository")
		}

		name := args[0]
		alias := args[1]
		repoBranches := internal.GetRepositoryBranches()
		branch, ok := repoBranches.GetBranchByName(name)
		if !ok {
			logger.Fatalln("Branch %v not found\n", name)
		}
		if repoBranches.AliasExists(alias) {
			logger.Fatalln("Alias already exists. Alias must be unique.")
		}
		if repoBranches.BranchWithAliasExists(alias) {
			logger.FatalF("A branch with name \"%s\" already exists. Alias must be unique.\n", alias)
		}

		branch.Alias = alias
		repoBranches.UpdateBranch(branch)
		logger.InfoF("Branch \"%v\" alias updated to \"%v\"\n", branch.Name, alias)
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}
