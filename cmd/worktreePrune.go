package cmd

import (
	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/cyrus2281/go-logger"
	"github.com/spf13/cobra"
)

var worktreePruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove stale worktree entries",
	Long: `Removes any worktrees that have been deleted outside of GitBranchTool.

First runs "git worktree prune" to clean up git's internal state,
then updates the internal list of worktrees to remove any entries
whose paths no longer appear in git's worktree list.

Example:
  g worktree prune`,
	Aliases: []string{"p"},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		git := internal.Git{}
		if !git.IsGitRepo() {
			logger.Fatalln("Not a git repository")
		}

		// First, run git worktree prune
		err := git.WorktreePrune()
		if err != nil {
			logger.Fatalln("Failed to prune git worktrees:", err)
		}

		// Get current worktree paths from git
		worktreeListOutput, err := git.WorktreeList()
		if err != nil {
			logger.Fatalln("Failed to list worktrees:", err)
		}
		worktreeMap := internal.ParseWorktreeList(worktreeListOutput)

		// Build set of active worktree paths
		activePaths := make(map[string]bool)
		for path := range worktreeMap {
			activePaths[path] = true
		}

		// Compare stored worktrees against active paths
		repoBranches := internal.GetRepositoryBranches()
		storedWorktrees := repoBranches.GetWorktrees()


		var toRemove []internal.Worktree
		for _, wt := range storedWorktrees {
			if !activePaths[wt.Path] {
				toRemove = append(toRemove, wt)
			}
		}
		for _, wt := range toRemove {
			repoBranches.RemoveWorktree(wt)
			logger.InfoF("Pruned stale worktree entry \"%s\" at %s\n", wt.Alias, wt.Path)
		}
		prunedCount := len(toRemove)

		if prunedCount == 0 {
			logger.Infoln("No stale worktree entries found")
		} else {
			logger.InfoF("Pruned %d stale worktree(s)\n", prunedCount)
		}
	},
}

func init() {
	worktreeCmd.AddCommand(worktreePruneCmd)
}
