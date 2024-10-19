/*
Copyright © 2024 Cyrus Mobini
*/
package cmd

import (
	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/cyrus2281/go-logger"
	"github.com/spf13/cobra"
)

// addAliasCmd represents the addAlias command
var addAliasCmd = &cobra.Command{
	Use:     "addAlias NAME ALIAS [...NOTE]",
	Short:   "Adds alias and note to a branch that is not stored yet",
	Long:    `Adds alias and note to a branch that is not stored yet`,
	Args:    cobra.MinimumNArgs(2),
	Aliases: []string{"a"},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return internal.GetGitBranches()
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		alias := args[1]
		notes := args[2:]
		notesString := ""
		for _, note := range notes {
			notesString += note + " "
		}

		newBranch := internal.Branch{
			Name:  name,
			Alias: alias,
			Note:  notesString,
		}

		git := internal.Git{}
		if !git.IsGitRepo() {
			logger.Fatalln("Not a git repository")
		}

		repoBranches := internal.GetRepositoryBranches()
		if repoBranches.AliasExists(newBranch.Alias) {
			logger.Fatalln("Alias already exists. Alias must be unique")
		}

		repoBranches.AddBranch(newBranch)
		logger.InfoF("Alias %v with note \"%v\" was added to branch %v", newBranch.Alias, newBranch.Note, newBranch.Name)
	},
}

func init() {
	rootCmd.AddCommand(addAliasCmd)
}
