/*
Copyright © 2024 Cyrus Mobini
*/
package cmd

import (
	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/cyrus2281/go-logger"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [...NAME|ALIAS]",
	Short: "Deletes listed branches base on name or alias",
	Long: `Deletes listed branches base on name or alias (requires at least one name/alias)"

	Without safe-delete uses the git command \"git branch -D [...NAME|ALIAS] \"
	With safe-delete uses the git command \"git branch [...NAME|ALIAS] \"`,
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"del", "d"},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return internal.GetBranchesAndAliases()
	},
	Run: func(cmd *cobra.Command, args []string) {
		git := internal.Git{}
		if !git.IsGitRepo() {
			logger.Fatalln("Not a git repository")
		}

		itemsToDelete := args
		safeDelete, _ := cmd.Flags().GetBool("safe-delete")
		ignoreErrors, _ := cmd.Flags().GetBool("ignore-errors")
		remote, _ := cmd.Flags().GetBool("remote")
		remoteOnly, _ := cmd.Flags().GetBool("remote-only")
		repoBranches := internal.GetRepositoryBranches()

		for _, item := range itemsToDelete {
			branch, ok := repoBranches.GetBranchByNameOrAlias(item)
			if !ok {
				logger.InfoF("Branch/Alias \"%v\" not found\n", item)
				continue
			}
			var err error
			if !remoteOnly {
				err = git.DeleteBranch(branch.Name, !safeDelete)
				if err != nil {
					logger.WarningF("Failed to delete branch \"%v\", %v\n", branch.Name, err)
				}

				if err == nil || ignoreErrors {
					repoBranches.RemoveBranch(branch)
					logger.InfoF("Branch \"%v\" with alias \"%v\" was deleted\n", branch.Name, branch.Alias)
				}
			}
			if (err == nil || ignoreErrors) && (remote || remoteOnly) {
				err = git.DeleteRemoteBranch(branch.Name)
				if err != nil {
					logger.WarningF("Failed to delete remote branch \"%v\", %v\n", branch.Name, err)
				} else {
					logger.InfoF("Remote branch \"%v\" was deleted\n", branch.Name)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolP("safe-delete", "s", false, "Safe delete branches - prevents deleting unmerged branches")
	deleteCmd.Flags().BoolP("ignore-errors", "i", false, "Ignore if git command fails and proceeds to remove the alias from the repository")
	deleteCmd.Flags().BoolP("remote", "r", false, "Delete the remote branch as well")
	deleteCmd.Flags().Bool("remote-only", false, "Deletes only the remote branch. Local branch and registry entry are not removed")
}
