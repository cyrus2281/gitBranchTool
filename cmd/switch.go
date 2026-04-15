package cmd

import (
	"strings"

	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/cyrus2281/go-logger"
	"github.com/spf13/cobra"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch NAME/ALIAS",
	Short: "Switches to the branch with the given name or alias",
	Long: `Switches to the branch with the given name or alias
Uses the git command \"git checkout NAME\"

This command can also be used to switch to and register a new branch at the same time.
For example:
	g switch NAME ALIAS [...NOTE]`,
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"checkout", "check", "s"},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		position := len(args) + 1
		if position == 1 {
			return internal.GetAllBranchesAndAliases()
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		git := internal.Git{}
		if !git.IsGitRepo() {
			logger.Fatalln("Not a git repository")
		}

		id := args[0]
		alias := ""
		notes := ""
		if len(args) > 1 {
			alias = args[1]
			if len(args) > 2 {
				for _, note := range args[2:] {
					notes += note + " "
				}
			}
		}

		repoBranches := internal.GetRepositoryBranches()

		worktreeAlias, _ := cmd.Flags().GetString("worktree")
		useWorktree := cmd.Flags().Changed("worktree")
		worktreeAlias = strings.TrimSpace(worktreeAlias)

		opts := switchOpts{
			UseWorktree:   useWorktree,
			WorktreeAlias: worktreeAlias,
		}

		if useWorktree {
			repoPath, err := git.GetRepositoryPath()
			if err != nil {
				logger.Fatalln("Failed to get repository path:", err)
			}
			repoName, err := git.GetRepositoryName()
			if err != nil {
				logger.Fatalln("Failed to get repository name:", err)
			}
			opts.WorktreePathTemplate = internal.GetWorktreePath()
			opts.RepoPath = repoPath
			opts.RepoName = repoName
		}

		executeSwitch(&git, &repoBranches, id, alias, notes, opts)
	},
}

type switchOpts struct {
	UseWorktree          bool
	WorktreeAlias        string
	WorktreePathTemplate string
	RepoPath             string
	RepoName             string
}

func executeSwitch(git *internal.Git, repoBranches *internal.RepositoryBranches,
	id, alias, notes string, opts switchOpts) {

	hasAlias := alias != ""
	branch, branchRegistered := repoBranches.GetBranchByNameOrAlias(id)

	if opts.UseWorktree {
		branchName := id
		if branchRegistered {
			branchName = branch.Name
		}

		worktreeListOutput, err := git.WorktreeList()
		if err != nil {
			logger.Fatalln("Failed to list worktrees:", err)
		}
		worktreeMap := internal.ParseWorktreeList(worktreeListOutput)
		existingPath := internal.GetWorktreePathForBranch(worktreeMap, branchName)

		if existingPath != "" {
			logger.InfoF("Branch \"%s\" checked out in worktree at %s\n", branchName, existingPath)
			return
		}

		wtAlias := opts.WorktreeAlias
		if wtAlias == "" || wtAlias == " " {
			if branchRegistered {
				wtAlias = branch.Alias
			} else if hasAlias {
				wtAlias = alias
			} else {
				wtAlias = branchName
			}
		}

		if !branchRegistered && hasAlias {
			newBranch := internal.Branch{Name: branchName, Alias: alias, Note: notes}
			if !repoBranches.AliasExists(alias) {
				repoBranches.AddBranch(newBranch)
				logger.InfoF("Branch \"%s\" has been registered with alias \"%s\"\n", branchName, alias)
			} else {
				logger.WarningF("Alias \"%s\" already exists. Branch was not registered.\n", alias)
			}
		}

		if repoBranches.WorktreeAliasExists(wtAlias) {
			logger.FatalF("Worktree alias \"%s\" already exists. Alias must be unique.\n", wtAlias)
		}

		resolvedPath := internal.ResolveWorktreePath(opts.WorktreePathTemplate, opts.RepoPath, opts.RepoName, wtAlias, branchName)

		err = git.WorktreeAdd(resolvedPath, branchName)
		if err != nil {
			logger.Fatalln("Failed to create worktree:", err)
		}

		newWorktree := internal.Worktree{Alias: wtAlias, Path: resolvedPath, Note: notes}
		repoBranches.AddWorktree(newWorktree)
		logger.InfoF("Worktree \"%s\" created at %s\n", wtAlias, resolvedPath)
		return
	}

	// Normal switch mode (no worktree)
	var checkoutErr error
	if !branchRegistered {
		checkoutErr = git.SwitchBranch(id)
	} else {
		checkoutErr = git.SwitchBranch(branch.Name)
	}
	if checkoutErr != nil {
		logger.ErrorF("Failed to switch branch to \"%v\"\n", id)
		logger.Fatalln(checkoutErr)
	}
	logger.InfoF("Switched to branch \"%v\"\n", id)

	defaultBranch := repoBranches.GetDefaultBranch()
	if defaultBranch == id {
		return
	}

	if hasAlias && !branchRegistered {
		newBranch := internal.Branch{Name: id, Alias: alias, Note: notes}
		if repoBranches.AliasExists(newBranch.Alias) {
			logger.WarningF("Alias \"%v\" already exists. Alias must be unique\n", newBranch.Alias)
			return
		}
		repoBranches.AddBranch(newBranch)
		logger.InfoF("Branch \"%v\" has been registered with alias  \"%v\"\n", newBranch.Name, newBranch.Alias)
	} else if !hasAlias && !branchRegistered {
		logger.WarningF("Branch \"%v\" is not registered with alias\n", id)
	}
}

func init() {
	rootCmd.AddCommand(switchCmd)
	switchCmd.Flags().StringP("worktree", "w", "", "Find or create a worktree for the branch. Use -w=ALIAS to specify a custom worktree alias.")
	switchCmd.Flags().Lookup("worktree").NoOptDefVal = " "
}
