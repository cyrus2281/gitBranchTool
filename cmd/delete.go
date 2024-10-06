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

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [...NAME|ALIAS]",
	Short: "Deletes listed branches base on name or alias",
	Long: `Deletes listed branches base on name or alias (requires at least one name/alias)"

	Without safe-delete uses the git command \"git branch -D [...NAME|ALIAS] \"
	With safe-delete uses the git command \"git branch [...NAME|ALIAS] \"`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		git := internal.Git{}
		if !git.IsGitRepo() {
			log.Fatalln("Not a git repository")
		}

		itemsToDelete := args
		safeDelete, _ := cmd.Flags().GetBool("safe-delete")
		ignoreErrors, _ := cmd.Flags().GetBool("ignore-errors")
		repoBranches := internal.GetRepositoryBranches()

		for _, item := range itemsToDelete {
			branch, ok := repoBranches.GetBranchByNameOrAlias(item)
			if !ok {
				log.Printf("Branch/Alias \"%v\" not found\n", item)
				continue
			}
			err := git.DeleteBranch(branch.Name, !safeDelete)
			if err != nil {
				fmt.Printf("Failed to delete branch \"%v\", %v", branch.Name, err)
			}

			if err == nil || ignoreErrors {
				repoBranches.RemoveBranch(branch)
				fmt.Printf("Branch \"%v\" with alias \"%v\" was deleted\n", branch.Name, branch.Alias)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	// Adding alias
	deleteCmd.Aliases = []string{"del", "d"}
	deleteCmd.Flags().BoolP("safe-delete", "s", false, "Safe delete branches - prevents deleting unmerged branches")
	deleteCmd.Flags().BoolP("ignore-errors", "i", false, "Ignore if git command fails and proceeds to remove the alias from the repository")
}
