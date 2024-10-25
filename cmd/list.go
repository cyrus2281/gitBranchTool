/*
Copyright © 2024 Cyrus Mobini
*/
package cmd

import (
	"fmt"

	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Lists all branches with their name, alias, and notes",
	Long:    `Lists all branches with their name, alias, and notes`,
	Aliases: []string{"ls", "l"},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		repoBranches := internal.GetRepositoryBranches()
		internal.PrintTableHeader()
		for index, branch := range repoBranches.GetBranches() {
			fmt.Printf("%d) %v\n", index, branch.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
