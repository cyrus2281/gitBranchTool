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
	Annotations: map[string]string{
		manualAnnotation: `List all registered branches with their name, alias, note, and any associated worktree.
Flags: -a/--all (also list local git branches that are not registered with g, shown with empty alias and note after the registered ones).`,
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		all, _ := cmd.Flags().GetBool("all")
		repoBranches := internal.GetRepositoryBranches()
		branches := repoBranches.GetBranches()

		git := internal.Git{}

		// With --all, append every local git branch that is not registered,
		// listed after the registered branches with empty alias and note.
		if all {
			branches = appendUnregisteredBranches(&git, branches)
		}

		// Check git worktree list for ANY worktrees (not just stored ones)
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

		// Determine if any listed branch has a worktree
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

// appendUnregisteredBranches returns a new slice containing every registered
// branch followed by each local git branch that is not already registered.
// Unregistered branches carry only their name (empty alias and note). If the
// git branch list cannot be read, the registered branches are returned as-is.
func appendUnregisteredBranches(git *internal.Git, registered []internal.Branch) []internal.Branch {
	result := make([]internal.Branch, len(registered))
	copy(result, registered)

	gitBranches, err := git.GetBranches()
	if err != nil {
		return result
	}

	registeredNames := make(map[string]bool, len(registered))
	for _, b := range registered {
		registeredNames[b.Name] = true
	}
	for _, name := range gitBranches {
		if !registeredNames[name] {
			result = append(result, internal.Branch{Name: name})
		}
	}
	return result
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("all", "a", false, "Also list local git branches that are not registered with g (shown with empty alias and note)")
}
