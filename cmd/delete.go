package cmd

import (
	"fmt"

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

		safeDelete, _ := cmd.Flags().GetBool("safe-delete")
		ignoreErrors, _ := cmd.Flags().GetBool("ignore-errors")
		remote, _ := cmd.Flags().GetBool("remote")
		remoteOnly, _ := cmd.Flags().GetBool("remote-only")
		forceWorktree, _ := cmd.Flags().GetBool("worktree")
		repoBranches := internal.GetRepositoryBranches()

		deleteBranchesWorktree := internal.GetConfig(internal.DELETE_BRANCHES_WORKTREE_KEY)
		var worktreeMap map[string]string
		if forceWorktree || deleteBranchesWorktree == "true" || deleteBranchesWorktree == "" || deleteBranchesWorktree == "null" {
			worktreeListOutput, err := git.WorktreeList()
			if err == nil {
				worktreeMap = internal.ParseWorktreeList(worktreeListOutput)
			}
		}

		for _, item := range args {
			shouldDeleteWt := false
			if worktreeMap != nil {
				branch, ok := repoBranches.GetBranchByNameOrAlias(item)
				if ok {
					wtPath := internal.GetWorktreePathForBranch(worktreeMap, branch.Name)
					if wtPath != "" {
						if forceWorktree || deleteBranchesWorktree == "true" {
							shouldDeleteWt = true
						} else if deleteBranchesWorktree != "false" {
							logger.InfoF("Branch \"%s\" is checked out in worktree at %s\n", branch.Name, wtPath)
							logger.InfoF("Delete the worktree as well? (y/n): ")
							var response string
							if _, err := fmt.Scanln(&response); err != nil {
								logger.WarningF("Failed to read response, defaulting to not deleting worktree: %v\n", err)
							} else {
								shouldDeleteWt = response == "y" || response == "Y" || response == "yes"
							}
						}
					}
				}
			}

			opts := deleteOpts{
				Force:                !safeDelete,
				IgnoreErrors:         ignoreErrors,
				Remote:               remote,
				RemoteOnly:           remoteOnly,
				ShouldDeleteWorktree: shouldDeleteWt,
				WorktreeMap:          worktreeMap,
			}
			executeDeleteBranch(&git, &repoBranches, item, opts)
		}
	},
}

type deleteOpts struct {
	Force                bool
	IgnoreErrors         bool
	Remote               bool
	RemoteOnly           bool
	ShouldDeleteWorktree bool
	WorktreeMap          map[string]string
}

func executeDeleteBranch(git *internal.Git, repoBranches *internal.RepositoryBranches,
	item string, opts deleteOpts) {

	branch, ok := repoBranches.GetBranchByNameOrAlias(item)
	if !ok {
		logger.InfoF("Branch/Alias \"%v\" not found\n", item)
		return
	}

	// Delete worktree FIRST if needed (git won't let us delete a branch checked out in a worktree)
	if opts.ShouldDeleteWorktree && opts.WorktreeMap != nil {
		worktreePath := internal.GetWorktreePathForBranch(opts.WorktreeMap, branch.Name)
		if worktreePath != "" {
			err := git.WorktreeRemove(worktreePath, true)
			if err != nil {
				logger.WarningF("Failed to delete worktree at %s: %v\n", worktreePath, err)
			} else {
				logger.InfoF("Worktree at %s was deleted\n", worktreePath)
				wt, found := repoBranches.GetWorktreeByPath(worktreePath)
				if found {
					repoBranches.RemoveWorktree(wt)
				}
			}
		}
	}

	var err error
	if !opts.RemoteOnly {
		err = git.DeleteBranch(branch.Name, opts.Force)
		if err != nil {
			logger.WarningF("Failed to delete branch \"%v\", %v\n", branch.Name, err)
		}

		if err == nil || opts.IgnoreErrors {
			repoBranches.RemoveBranch(branch)
			logger.InfoF("Branch \"%v\" with alias \"%v\" was deleted\n", branch.Name, branch.Alias)
		}
	}
	if (err == nil || opts.IgnoreErrors) && (opts.Remote || opts.RemoteOnly) {
		err = git.DeleteRemoteBranch(branch.Name)
		if err != nil {
			logger.WarningF("Failed to delete remote branch \"%v\", %v\n", branch.Name, err)
		} else {
			logger.InfoF("Remote branch \"%v\" was deleted\n", branch.Name)
		}
	}
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolP("safe-delete", "s", false, "Safe delete branches - prevents deleting unmerged branches")
	deleteCmd.Flags().BoolP("ignore-errors", "i", false, "Ignore if git command fails and proceeds to remove the alias from the repository")
	deleteCmd.Flags().BoolP("remote", "r", false, "Delete the remote branch as well")
	deleteCmd.Flags().Bool("remote-only", false, "Deletes only the remote branch. Local branch and registry entry are not removed")
	deleteCmd.Flags().BoolP("worktree", "w", false, "Also delete the associated worktree")
}
