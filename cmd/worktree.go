package cmd

import (
	"fmt"

	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/cyrus2281/go-logger"
	"github.com/spf13/cobra"
)

var worktreeCmd = &cobra.Command{
	Use:   "worktree",
	Short: "Manage git worktrees (Run `g worktree --help` for more information)",
	Long: `Manage git worktrees

- worktree create ALIAS [BRANCH] [...NOTE] : Create a new worktree
	Example: g worktree create my-feature feature/my-feature "Working on login page"
	Example: g worktree create quick-fix  (creates new branch "quick-fix")

- worktree list : List all worktrees
	Example: g worktree list

- worktree delete [...PATH|ALIAS] [--force|-f] : Delete worktree(s)
	Example: g worktree delete my-feature
	Example: g worktree delete my-feature --force

- worktree prune : Remove stale worktree entries
	Example: g worktree prune
`,
	Aliases: []string{"w"},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"create", "list", "delete", "prune"}, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		git := internal.Git{}
		if !git.IsGitRepo() {
			logger.Fatalln("Not a git repository")
		}

		// Print usage help
		cmd.Help()
		fmt.Println()

		// Run the worktree list
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
		shownPaths := make(map[string]bool)

		internal.PrintWorktreeTableHeader()
		index := 0

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
			unregistered := internal.Worktree{Path: path}
			logger.InfoF("%d) ", index)
			fmt.Println(unregistered.StringWithBranch(branch, branchAlias))
			index++
		}
	},
}

func init() {
	rootCmd.AddCommand(worktreeCmd)
}
