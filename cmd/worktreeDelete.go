package cmd

import (
	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/cyrus2281/go-logger"
	"github.com/spf13/cobra"
)

var worktreeDeleteCmd = &cobra.Command{
	Use:   "delete [...PATH|ALIAS]",
	Short: "Delete worktree(s) by path or alias",
	Long: `Deletes the specified worktree(s) by path or alias.

If the worktree has uncommitted changes, the deletion will fail unless --force is used.

Examples:
  g worktree delete my-feature
  g worktree delete my-feature bugfix-123
  g w d my-feature --force`,
	Aliases: []string{"d"},
	Args:    cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return internal.GetWorktreeAliasesAndPaths()
	},
	Run: func(cmd *cobra.Command, args []string) {
		git := internal.Git{}
		if !git.IsGitRepo() {
			logger.Fatalln("Not a git repository")
		}

		force, _ := cmd.Flags().GetBool("force")
		repoBranches := internal.GetRepositoryBranches()

		for _, item := range args {
			// Look up worktree by alias first, then by stored path
			wt, ok := repoBranches.GetWorktreeByAlias(item)
			if !ok {
				wt, ok = repoBranches.GetWorktreeByPath(item)
			}

			if ok {
				// Registered worktree — remove via stored path
				err := git.WorktreeRemove(wt.Path, force)
				if err != nil {
					logger.WarningF("Failed to delete worktree \"%s\": %v\n", wt.Alias, err)
					continue
				}
				repoBranches.RemoveWorktree(wt)
				logger.InfoF("Worktree \"%s\" at %s was deleted\n", wt.Alias, wt.Path)
			} else {
				// Not registered — try removing by path directly via git
				err := git.WorktreeRemove(item, force)
				if err != nil {
					logger.WarningF("Failed to delete worktree \"%s\": %v\n", item, err)
					continue
				}
				logger.InfoF("Worktree at %s was deleted\n", item)
			}
		}
	},
}

func init() {
	worktreeCmd.AddCommand(worktreeDeleteCmd)
	worktreeDeleteCmd.Flags().BoolP("force", "f", false, "Force delete worktree even with uncommitted changes")
}
