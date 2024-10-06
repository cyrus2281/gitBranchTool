/*
Copyright © 2024 Cyrus Mobini
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/spf13/cobra"
)

// addAliasCmd represents the addAlias command
var addAliasCmd = &cobra.Command{
	Use:     "addAlias NAME ALIAS [...NOTE]",
	Short:   "Adds alias and note to a branch that is not stored yet",
	Long:    `Adds alias and note to a branch that is not stored yet`,
	Args:    cobra.MinimumNArgs(2),
	Aliases: []string{"a"},
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
			log.Fatalln("Not a git repository")
		}

		repoBranches := internal.GetRepositoryBranches()
		if repoBranches.AliasExists(newBranch.Alias) {
			log.Fatalln("Alias already exists. Alias must be unique")
		}

		repoBranches.AddBranch(newBranch)
		fmt.Printf("Alias %v with note \"%v\" was added to branch %v\n", newBranch.Alias, newBranch.Note, newBranch.Name)
	},
}

func init() {
	rootCmd.AddCommand(addAliasCmd)
}
