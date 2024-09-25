/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all branches with their id, alias, and notes",
	Long:  `Lists all branches with their id, alias, and notes`,
	Run: func(cmd *cobra.Command, args []string) {
		git := internal.Git{}
		// Get the home directory
		gHome := viper.GetString("GIT_BRANCH_TOOL_HOME")
		if gHome == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				panic(err)
			}
			gHome = filepath.Join(home, ".gitBranchTool_go")
		}
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
