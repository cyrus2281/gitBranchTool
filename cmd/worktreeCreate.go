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
If no registered reference is found, it will create the worktree with the branch as-is.
If the branch does not exist in git, a new branch will be created but will NOT be
automatically registered (use "g addAlias <branch> <alias>" to register it).
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
		branchArg := ""
		notes := ""
		if len(args) > 1 {
			branchArg = args[1]
			if len(args) > 2 {
				for _, note := range args[2:] {
					notes += note + " "
				}
			}
		}

		repoBranches := internal.GetRepositoryBranches()

		repoPath, err := git.GetRepositoryPath()
		if err != nil {
			logger.Fatalln("Failed to get repository path:", err)
		}
		repoName, err := git.GetRepositoryName()
		if err != nil {
			logger.Fatalln("Failed to get repository name:", err)
		}
		template := internal.GetWorktreePath()

		executeWorktreeCreate(&git, &repoBranches, alias, branchArg, notes, template, repoPath, repoName)
	},
}

func executeWorktreeCreate(git *internal.Git, repoBranches *internal.RepositoryBranches,
	alias, branchArg, notes, worktreePathTemplate, repoPath, repoName string) {

	hasBranch := branchArg != ""

	if repoBranches.WorktreeAliasExists(alias) {
		logger.FatalF("Worktree alias \"%s\" already exists. Alias must be unique.\n", alias)
	}

	// Resolve the branch name and determine if we should auto-register
	var branchName string
	branchResolved := false
	if hasBranch {
		branch, ok := repoBranches.GetBranchByNameOrAlias(branchArg)
		if ok {
			branchName = branch.Name
			branchResolved = true
		} else {
			branchName = branchArg
		}
	} else {
		branchName = alias
	}

	shouldAutoRegister := !hasBranch || branchResolved
	if shouldAutoRegister {
		branchAlreadyRegistered := repoBranches.BranchExists(internal.Branch{Name: branchName})
		if !branchAlreadyRegistered && repoBranches.AliasExists(alias) {
			logger.FatalF("Branch alias \"%s\" already exists. Cannot auto-register branch \"%s\".\n", alias, branchName)
		}
	}

	resolvedPath := internal.ResolveWorktreePath(worktreePathTemplate, repoPath, repoName, alias, branchName)

	// Create the worktree
	var err error
	if !hasBranch {
		err = git.WorktreeAddNewBranch(branchName, resolvedPath)
	} else {
		err = git.WorktreeAdd(resolvedPath, branchName)
		if err != nil && !branchResolved {
			err = git.WorktreeAddNewBranch(branchName, resolvedPath)
		}
	}
	if err != nil {
		logger.Fatalln("Failed to create worktree:", err)
	}

	// Store the worktree
	newWorktree := internal.Worktree{Alias: alias, Path: resolvedPath, Note: notes}
	repoBranches.AddWorktree(newWorktree)
	logger.InfoF("Worktree \"%s\" created at %s\n", alias, resolvedPath)

	// Auto-register branch only when appropriate
	if shouldAutoRegister {
		branchAlreadyRegistered := repoBranches.BranchExists(internal.Branch{Name: branchName})
		if !branchAlreadyRegistered {
			branchEntry := internal.Branch{Name: branchName, Alias: alias, Note: notes}
			repoBranches.AddBranch(branchEntry)
			logger.InfoF("Branch \"%s\" registered with alias \"%s\"\n", branchName, alias)
		} else {
			logger.DebugF("Branch \"%s\" is already registered\n", branchName)
		}
	}
}

func init() {
	worktreeCmd.AddCommand(worktreeCreateCmd)
}
