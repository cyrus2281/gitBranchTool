/*
Copyright © 2024 Cyrus Mobini
*/
package cmd

import (
	"log"

	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/spf13/cobra"
)

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:     "rename NAME ALIAS",
	Short:   "Updates the alias for the given branch name",
	Long:    `Updates the alias for the given branch name.`,
	Args:    cobra.ExactArgs(2),
	Aliases: []string{"updateBranchAlias", "update-branch-alias", "uba"},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return internal.GetBranches()
	},
	Run: func(cmd *cobra.Command, args []string) {
		git := internal.Git{}
		if !git.IsGitRepo() {
			log.Fatalln("Not a git repository")
		}

		name := args[0]
		alias := args[1]
		repoBranches := internal.GetRepositoryBranches()
		branch, ok := repoBranches.GetBranchByName(name)
		if !ok {
			log.Fatalf("Branch %v not found\n", alias)
		}
		branch.Alias = alias
		repoBranches.UpdateBranch(branch)
		log.Printf("Branch \"%v\" alias updated to \"%v\"\n", branch.Name, alias)
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}
