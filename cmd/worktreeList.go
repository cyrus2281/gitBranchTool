package cmd

import (
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
		shownPaths := make(map[string]bool)

		headers := []string{"Path", "Alias", "Branch", "Branch Alias", "Note"}
		rows := [][]string{}

		// First add stored worktrees (these have alias/note)
		for _, wt := range storedWorktrees {
			branch := worktreeMap[wt.Path]
			branchAlias := ""
			if branch != "" {
				b, ok := repoBranches.GetBranchByName(branch)
				if ok {
					branchAlias = b.Alias
				}
			}
			rows = append(rows, []string{wt.Path, wt.Alias, branch, branchAlias, wt.Note})
			shownPaths[wt.Path] = true
		}

		// Then add unregistered worktrees from git
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
			rows = append(rows, []string{path, "", branch, branchAlias, ""})
		}

		internal.PrintTable(headers, rows)
	},
}

func init() {
	worktreeCmd.AddCommand(worktreeListCmd)
}
