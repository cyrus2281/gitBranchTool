/*
Copyright © 2024 Cyrus Mobini
*/
package cmd

import (
	"log"

	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/spf13/cobra"
)

// setDefaultBranchCmd represents the setDefaultBranch command
var setDefaultBranchCmd = &cobra.Command{
	Use:     "setDefaultBranch NAME",
	Short:   "Change the default branch, default is " + internal.DEFAULT_BRANCH,
	Long:    "Change the default branch, default is " + internal.DEFAULT_BRANCH,
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"set-default-branch", "sdb"},
	Run: func(cmd *cobra.Command, args []string) {
		git := internal.Git{}
		if !git.IsGitRepo() {
			log.Fatalln("Not a git repository")
		}
		name := args[0]
		repoBranches := internal.GetRepositoryBranches()
		repoBranches.SetDefaultBranch(name)
		log.Printf("Default branch set to \"%v\"\n", name)
	},
}

func init() {
	rootCmd.AddCommand(setDefaultBranchCmd)
}
