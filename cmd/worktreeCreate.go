package cmd

import (
	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/cyrus2281/go-logger"
	"github.com/spf13/cobra"
)

var worktreeCreateCmd = &cobra.Command{
	Use:   "create ALIAS [BRANCH] [...NOTE]",
	Short: "Create a new worktree",
	Long: `Create a new worktree for a branch.

If BRANCH is provided, it can be a branch name or alias (resolved from registered branches).
If BRANCH is not provided, a new branch is created with the alias as the branch name.

The worktree path is determined by the worktree-path setting.
Use "g set worktree-path" to customize the path template.
Available variables: {repository}, {alias}, {branch}

Examples:
  g worktree create my-feature feature/my-feature "Working on login page"
  g w c my-feature mf
  g worktree create quick-fix`,
	Aliases: []string{"c"},
	Args:    cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		position := len(args) + 1
		if position == 2 {
			return internal.GetAllBranchesAndAliases()
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		git := internal.Git{}
		if !git.IsGitRepo() {
			logger.Fatalln("Not a git repository")
		}

		alias := args[0]
		hasBranch := len(args) > 1
		branchArg := ""
		notes := ""

		if hasBranch {
			branchArg = args[1]
			if len(args) > 2 {
				for _, note := range args[2:] {
					notes += note + " "
				}
			}
		}

		repoBranches := internal.GetRepositoryBranches()

		// Check if worktree alias already exists
		if repoBranches.WorktreeAliasExists(alias) {
			logger.FatalF("Worktree alias \"%s\" already exists. Alias must be unique.\n", alias)
		}

		// Resolve the branch name
		var branchName string
		if hasBranch {
			// Try to resolve from registered branches
			branch, ok := repoBranches.GetBranchByNameOrAlias(branchArg)
			if ok {
				branchName = branch.Name
			} else {
				// Use as-is (could be an unregistered git branch)
				branchName = branchArg
			}
		} else {
			// No branch provided — will create a new branch with alias as name
			branchName = alias
		}

		// Resolve worktree path
		repoPath, err := git.GetRepositoryPath()
		if err != nil {
			logger.Fatalln("Failed to get repository path:", err)
		}
		repoName, err := git.GetRepositoryName()
		if err != nil {
			logger.Fatalln("Failed to get repository name:", err)
		}
		template := internal.GetWorktreePath()
		resolvedPath := internal.ResolveWorktreePath(template, repoPath, repoName, alias, branchName)

		// Create the worktree
		if hasBranch {
			// Worktree for existing branch
			err = git.WorktreeAdd(resolvedPath, branchName)
		} else {
			// New branch + worktree
			err = git.WorktreeAddNewBranch(branchName, resolvedPath)
		}
		if err != nil {
			logger.Fatalln("Failed to create worktree:", err)
		}

		// Store the worktree
		newWorktree := internal.Worktree{
			Alias: alias,
			Path:  resolvedPath,
			Note:  notes,
		}
		repoBranches.AddWorktree(newWorktree)
		logger.InfoF("Worktree \"%s\" created at %s\n", alias, resolvedPath)

		// Auto-register branch if not already registered
		branchEntry := internal.Branch{
			Name:  branchName,
			Alias: alias,
			Note:  notes,
		}
		if repoBranches.BranchExists(branchEntry) {
			logger.DebugF("Branch \"%s\" is already registered\n", branchName)
		} else {
			if repoBranches.AliasExists(alias) {
				logger.WarningF("Branch alias \"%s\" already exists. Branch \"%s\" was not auto-registered.\n", alias, branchName)
			} else {
				repoBranches.AddBranch(branchEntry)
				logger.InfoF("Branch \"%s\" registered with alias \"%s\"\n", branchName, alias)
			}
		}
	},
}

func init() {
	worktreeCmd.AddCommand(worktreeCreateCmd)
}
