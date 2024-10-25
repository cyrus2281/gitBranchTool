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

// getBranchAliasCmd represents the getBranchAlias command
var getBranchAliasCmd = &cobra.Command{
	Use:     "getBranchAlias NAME",
	Short:   "Gets the branch alias",
	Long:    `Gets the branch alias`,
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"get-branch-alias", "gbr"},
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
		repoBranches := internal.GetRepositoryBranches()
		branch, ok := repoBranches.GetBranchByName(name)
		if !ok {
			logger.FatalF("Branch %v not found\n", name)
		}

		fmt.Println(branch.Alias)
	},
}

func init() {
	rootCmd.AddCommand(getBranchAliasCmd)
}
