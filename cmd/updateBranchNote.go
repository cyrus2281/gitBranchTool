/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/spf13/cobra"
)

// updateBranchNoteCmd represents the updateBranchNote command
var updateBranchNoteCmd = &cobra.Command{
	Use:   "updateBranchNote NAME/ALIAS [...NOTE]",
	Short: "Adds/updates the notes for a branch base on name/alias",
	Long:  `Adds/updates the notes for a branch base on name/alias`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		git := internal.Git{}
		if !git.IsGitRepo() {
			log.Fatalln("Not a git repository")
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
			log.Fatalf("Branch %v not found\n", id)
		}
		branch.Note = notesString
		repoBranches.UpdateBranch(branch)
		log.Printf("Branch \"%v\" notes were updated\n", branch.Name)
	},
}

func init() {
	rootCmd.AddCommand(updateBranchNoteCmd)
	updateBranchNoteCmd.Aliases = []string{"update-branch-note", "ubn"}
}
