package cmd

import (
	"fmt"

	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/cyrus2281/go-logger"
	"github.com/spf13/cobra"
)

var worktreeListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all worktrees",
	Long:    `Lists all worktrees associated with the repository, showing their path, alias, branch, branch alias, and note.`,
	Aliases: []string{"l"},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		git := internal.Git{}
		if !git.IsGitRepo() {
			logger.Fatalln("Not a git repository")
		}

		worktreeListOutput, err := git.WorktreeList()
		if err != nil {
			logger.Fatalln("Failed to list worktrees:", err)
		}
		worktreeMap := internal.ParseWorktreeList(worktreeListOutput)

		if len(worktreeMap) == 0 {
			logger.Infoln("No worktrees found")
			return
		}

		repoBranches := internal.GetRepositoryBranches()
		storedWorktrees := repoBranches.GetWorktrees()

		// Build a set of paths we've already shown from stored worktrees
		shownPaths := make(map[string]bool)

		internal.PrintWorktreeTableHeader()
		index := 0

		// First show stored worktrees (these have alias/note)
		for _, wt := range storedWorktrees {
			branch := worktreeMap[wt.Path]
			branchAlias := ""
			if branch != "" {
				b, ok := repoBranches.GetBranchByName(branch)
				if ok {
					branchAlias = b.Alias
				}
			}
			logger.InfoF("%d) ", index)
			fmt.Println(wt.StringWithBranch(branch, branchAlias))
			shownPaths[wt.Path] = true
			index++
		}

		// Then show unregistered worktrees from git
		for path, branch := range worktreeMap {
			if shownPaths[path] {
				continue
			}
			branchAlias := ""
			if branch != "" {
				b, ok := repoBranches.GetBranchByName(branch)
				if ok {
					branchAlias = b.Alias
				}
			}
			unregistered := internal.Worktree{
				Alias: "",
				Path:  path,
				Note:  "",
			}
			logger.InfoF("%d) ", index)
			fmt.Println(unregistered.StringWithBranch(branch, branchAlias))
			index++
		}
	},
}

func init() {
	worktreeCmd.AddCommand(worktreeListCmd)
}
