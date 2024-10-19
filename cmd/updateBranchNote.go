/*
Copyright © 2024 Cyrus Mobini
*/
package cmd

import (
	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/cyrus2281/go-logger"
	"github.com/spf13/cobra"
)

// updateBranchNoteCmd represents the updateBranchNote command
var updateBranchNoteCmd = &cobra.Command{
	Use:   "updateBranchNote NAME/ALIAS [...NOTE]",
	Short: "Adds/updates the notes for a branch base on name/alias",
	Long:  `Adds/updates the notes for a branch base on name/alias`,
	Args:  cobra.MinimumNArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return internal.GetBranchesAndAliases()
	},
	Aliases: []string{"update-branch-note", "ubn"},
	Run: func(cmd *cobra.Command, args []string) {
		git := internal.Git{}
		if !git.IsGitRepo() {
			logger.Fatalln("Not a git repository")
		}

		id := args[0]
		notes := args[1:]
		notesString := ""
		for _, note := range notes {
			notesString += note + " "
		}
		repoBranches := internal.GetRepositoryBranches()
		branch, ok := repoBranches.GetBranchByNameOrAlias(id)
		if !ok {
			logger.FatalF("Branch %v not found\n", id)
		}
		branch.Note = notesString
		repoBranches.UpdateBranch(branch)
		logger.InfoF("Branch \"%v\" notes were updated\n", branch.Name)
	},
}

func init() {
	rootCmd.AddCommand(updateBranchNoteCmd)
}
