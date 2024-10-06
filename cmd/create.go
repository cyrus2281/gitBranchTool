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

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create NAME ALIAS [...NOTE]",
	Short: "Creates a branch with name, alias, and note, and checks into it",
	Long: `Creates a branch with name, alias, and note, and checks into it.
	Without only-create flag: "git checkout -b NAME"
	With only-create flag: "git branch NAME"
	`,
	Aliases: []string{"c"},
	Args:    cobra.MinimumNArgs(2),
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

		createOnly, _ := cmd.Flags().GetBool("only-create")
		var err error
		if createOnly {
			err = git.CreateNewBranch(newBranch.Name)
		} else {
			err = git.SwitchToNewBranch(newBranch.Name)
		}
		if err != nil {
			log.Fatalln(err)
		}

		repoBranches.AddBranch(newBranch)
		fmt.Printf("Branch %v with alias %v was created\n", newBranch.Name, newBranch.Alias)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().BoolP("only-create", "o", false, "Only create the branch, do not check into it")
}
