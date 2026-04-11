package cmd

import (
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
		branches := repoBranches.GetBranches()

		// Check git worktree list for ANY worktrees (not just stored ones)
		git := internal.Git{}
		worktreeListOutput, _ := git.WorktreeList()
		worktreeMap := internal.ParseWorktreeList(worktreeListOutput)

		// Build branch-to-worktree info map from git worktree list
		// Map branch name -> stored alias if available, otherwise the worktree path
		branchToWorktreeInfo := make(map[string]string)
		for path, branchName := range worktreeMap {
			if branchName == "" {
				continue
			}
			wt, ok := repoBranches.GetWorktreeByPath(path)
			if ok {
				branchToWorktreeInfo[branchName] = wt.Alias
			} else {
				branchToWorktreeInfo[branchName] = path
			}
		}

		// Determine if any registered branch has a worktree
		hasWorktrees := false
		for _, branch := range branches {
			if _, ok := branchToWorktreeInfo[branch.Name]; ok {
				hasWorktrees = true
				break
			}
		}

		if hasWorktrees {
			headers := []string{"Branch Name", "Alias", "Note", "Worktree"}
			rows := make([][]string, len(branches))
			for i, branch := range branches {
				rows[i] = []string{branch.Name, branch.Alias, branch.Note, branchToWorktreeInfo[branch.Name]}
			}
			internal.PrintTable(headers, rows)
		} else {
			headers := []string{"Branch Name", "Alias", "Note"}
			rows := make([][]string, len(branches))
			for i, branch := range branches {
				rows[i] = []string{branch.Name, branch.Alias, branch.Note}
			}
			internal.PrintTable(headers, rows)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
