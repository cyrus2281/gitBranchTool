/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all branches with their id, alias, and notes",
	Long:  `Lists all branches with their id, alias, and notes`,
	Run: func(cmd *cobra.Command, args []string) {
		git := internal.Git{}
		// Get the home directory
		gHome := internal.GetHome()
		repositoryName, err := git.GetRepositoryName()
		if err != nil {
			panic(err)
		}
		repoBranches := internal.RepositoryBranches{
			RepositoryName: repositoryName,
			StoreDirectory: gHome,
		}
		for index, branch := range repoBranches.GetBranches() {
			println(fmt.Sprintf("%d) %v", index, branch.Print()))
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	// Adding alias
	listCmd.Aliases = []string{"ls", "l"}
}
