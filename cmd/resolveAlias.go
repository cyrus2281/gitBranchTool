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

// resolveAliasCmd represents the resolveAlias command
var resolveAliasCmd = &cobra.Command{
	Use:   "resolveAlias ALIAS",
	Short: "Resolves the branch name from an alias",
	Long: `Resolves the branch name from an alias
	
	Example: git merge $(g r ALIAS)`,
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"resolve-alias", "r"},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return internal.GetAliases()
	},
	Run: func(cmd *cobra.Command, args []string) {
		git := internal.Git{}
		if !git.IsGitRepo() {
			logger.Fatalln("Not a git repository")
		}

		alias := args[0]
		repoBranches := internal.GetRepositoryBranches()
		branch, ok := repoBranches.GetBranchByAlias(alias)
		if !ok {
			logger.FatalF("Alias %v not found\n", alias)
		}

		fmt.Println(branch.Name)
	},
}

func init() {
	rootCmd.AddCommand(resolveAliasCmd)
}
