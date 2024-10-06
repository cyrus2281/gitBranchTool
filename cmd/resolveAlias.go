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

// resolveAliasCmd represents the resolveAlias command
var resolveAliasCmd = &cobra.Command{
	Use:   "resolveAlias ALIAS",
	Short: "Resolves the branch name from an alias",
	Long: `Resolves the branch name from an alias
	
	Example: git merge $(g r ALIAS)`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		git := internal.Git{}
		if !git.IsGitRepo() {
			log.Fatalln("Not a git repository")
		}

		alias := args[0]
		repoBranches := internal.GetRepositoryBranches()
		branch, ok := repoBranches.GetBranchByAlias(alias)
		if !ok {
			log.Fatalf("Alias %v not found\n", alias)
		}

		fmt.Println(branch.Name)
	},
}

func init() {
	rootCmd.AddCommand(resolveAliasCmd)
	resolveAliasCmd.Aliases = []string{"r"}

}
