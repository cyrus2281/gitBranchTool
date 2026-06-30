package cmd

import (
	"fmt"
	"regexp"

	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/cyrus2281/go-logger"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [...NAME|ALIAS]",
	Short: "Deletes listed branches base on name or alias",
	Long: `Deletes listed branches base on name or alias (requires at least one name/alias)"

	Without safe-delete uses the git command \"git branch -D [...NAME|ALIAS] \"
	With safe-delete uses the git command \"git branch [...NAME|ALIAS] \"

	With --regex/-e, the arguments are treated as regular expressions and every
	registered branch whose name or alias matches is deleted.

	With --all/-a, branches that are not registered with g are also considered:
	the candidate list is expanded with every local git branch before the
	search runs, so unregistered branches can be deleted by name or by regex.`,
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"del", "d"},
	Annotations: map[string]string{
		manualAnnotation: `Delete one or more registered branches by NAME or ALIAS, removing both the git branch and its registry entry.
Flags: -s/--safe-delete (refuse unmerged branches), -r/--remote (also delete the remote branch), --remote-only, -i/--ignore-errors, -w/--worktree (also remove its worktree), -e/--regex (treat arguments as regular expressions matched against branch names and aliases), -a/--all (also consider local git branches not registered with g, expanding the candidate list before searching/deleting).`,
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if all, _ := cmd.Flags().GetBool("all"); all {
			return internal.GetAllBranchesAndAliases()
		}
		return internal.GetBranchesAndAliases()
	},
	Run: func(cmd *cobra.Command, args []string) {
		git := internal.Git{}
		if !git.IsGitRepo() {
			logger.Fatalln("Not a git repository")
		}

		safeDelete, _ := cmd.Flags().GetBool("safe-delete")
		ignoreErrors, _ := cmd.Flags().GetBool("ignore-errors")
		remote, _ := cmd.Flags().GetBool("remote")
		remoteOnly, _ := cmd.Flags().GetBool("remote-only")
		forceWorktree, _ := cmd.Flags().GetBool("worktree")
		regex, _ := cmd.Flags().GetBool("regex")
		all, _ := cmd.Flags().GetBool("all")
		repoBranches := internal.GetRepositoryBranches()

		// With --all, unregistered git branches are also candidates for deletion.
		var gitBranches []string
		if all {
			gitBranches, _ = git.GetBranches()
		}

		targets, err := resolveDeleteTargets(&repoBranches, gitBranches, args, regex, all)
		if err != nil {
			logger.Fatalln(err)
		}
		if regex && len(targets) == 0 {
			logger.InfoF("No branches matched the given pattern(s)\n")
			return
		}

		deleteBranchesWorktree := internal.GetConfig(internal.DELETE_BRANCHES_WORKTREE_KEY)
		var worktreeMap map[string]string
		if forceWorktree || deleteBranchesWorktree == "true" || deleteBranchesWorktree == "" || deleteBranchesWorktree == "null" {
			worktreeListOutput, err := git.WorktreeList()
			if err == nil {
				worktreeMap = internal.ParseWorktreeList(worktreeListOutput)
			}
		}

		for _, item := range targets {
			shouldDeleteWt := false
			if worktreeMap != nil {
				// Resolve to the branch name, falling back to the item itself
				// so unregistered branches (only reachable with --all) still
				// have their worktree detected.
				branchName := item
				if branch, ok := repoBranches.GetBranchByNameOrAlias(item); ok {
					branchName = branch.Name
				}
				wtPath := internal.GetWorktreePathForBranch(worktreeMap, branchName)
				if wtPath != "" {
					if forceWorktree || deleteBranchesWorktree == "true" {
						shouldDeleteWt = true
					} else if deleteBranchesWorktree != "false" {
						logger.InfoF("Branch \"%s\" is checked out in worktree at %s\n", branchName, wtPath)
						logger.InfoF("Delete the worktree as well? (y/n): ")
						var response string
						if _, err := fmt.Scanln(&response); err != nil {
							logger.WarningF("Failed to read response, defaulting to not deleting worktree: %v\n", err)
						} else {
							shouldDeleteWt = response == "y" || response == "Y" || response == "yes"
						}
					}
				}
			}

			opts := deleteOpts{
				Force:                !safeDelete,
				IgnoreErrors:         ignoreErrors,
				Remote:               remote,
				RemoteOnly:           remoteOnly,
				ShouldDeleteWorktree: shouldDeleteWt,
				WorktreeMap:          worktreeMap,
				All:                  all,
			}
			executeDeleteBranch(&git, &repoBranches, item, opts)
		}
	},
}

type deleteOpts struct {
	Force                bool
	IgnoreErrors         bool
	Remote               bool
	RemoteOnly           bool
	ShouldDeleteWorktree bool
	WorktreeMap          map[string]string
	// All allows deleting a git branch that is not registered with g, treating
	// the item as a raw branch name when it is not found in the registry.
	All bool
}

// resolveDeleteTargets expands the provided args into the list of items to delete.
// When useRegex is false, the args are returned unchanged (each treated as a
// literal NAME or ALIAS). When useRegex is true, each arg is compiled as a
// regular expression and matched against every registered branch's name and
// alias; the names of all matching branches are returned (de-duplicated, in
// registry order).
//
// When all is true, the regex search is additionally run against gitBranches
// (the local git branch names), so unregistered branches matched by a pattern
// are appended after the registered matches. all has no effect when useRegex
// is false, where unregistered names are handled at delete time instead.
func resolveDeleteTargets(repoBranches *internal.RepositoryBranches, gitBranches []string, args []string, useRegex bool, all bool) ([]string, error) {
	if !useRegex {
		return args, nil
	}

	patterns := make([]*regexp.Regexp, 0, len(args))
	for _, arg := range args {
		re, err := regexp.Compile(arg)
		if err != nil {
			return nil, fmt.Errorf("invalid regular expression \"%v\": %w", arg, err)
		}
		patterns = append(patterns, re)
	}

	seen := make(map[string]bool)
	targets := []string{}
	for _, branch := range repoBranches.GetBranches() {
		for _, re := range patterns {
			if re.MatchString(branch.Name) || re.MatchString(branch.Alias) {
				if !seen[branch.Name] {
					seen[branch.Name] = true
					targets = append(targets, branch.Name)
				}
				break
			}
		}
	}
	if all {
		for _, name := range gitBranches {
			if seen[name] {
				continue
			}
			for _, re := range patterns {
				if re.MatchString(name) {
					seen[name] = true
					targets = append(targets, name)
					break
				}
			}
		}
	}
	return targets, nil
}

func executeDeleteBranch(git *internal.Git, repoBranches *internal.RepositoryBranches,
	item string, opts deleteOpts) {

	branch, ok := repoBranches.GetBranchByNameOrAlias(item)
	if !ok {
		if !opts.All {
			logger.InfoF("Branch/Alias \"%v\" not found\n", item)
			return
		}
		// With --all, the item is treated as a raw git branch name. There is no
		// registry entry to remove; the git delete below reports an error if no
		// such branch exists.
		branch = internal.Branch{Name: item}
	}

	// Delete worktree FIRST if needed (git won't let us delete a branch checked out in a worktree)
	if opts.ShouldDeleteWorktree && opts.WorktreeMap != nil {
		worktreePath := internal.GetWorktreePathForBranch(opts.WorktreeMap, branch.Name)
		if worktreePath != "" {
			err := git.WorktreeRemove(worktreePath, true)
			if err != nil {
				logger.WarningF("Failed to delete worktree at %s: %v\n", worktreePath, err)
			} else {
				logger.InfoF("Worktree at %s was deleted\n", worktreePath)
				wt, found := repoBranches.GetWorktreeByPath(worktreePath)
				if found {
					repoBranches.RemoveWorktree(wt)
				}
			}
		}
	}

	var err error
	if !opts.RemoteOnly {
		err = git.DeleteBranch(branch.Name, opts.Force)
		if err != nil {
			logger.WarningF("Failed to delete branch \"%v\", %v\n", branch.Name, err)
		}

		if err == nil || opts.IgnoreErrors {
			// RemoveBranch is a no-op for unregistered branches (no matching entry).
			repoBranches.RemoveBranch(branch)
			if branch.Alias != "" {
				logger.InfoF("Branch \"%v\" with alias \"%v\" was deleted\n", branch.Name, branch.Alias)
			} else {
				logger.InfoF("Branch \"%v\" was deleted\n", branch.Name)
			}
		}
	}
	if (err == nil || opts.IgnoreErrors) && (opts.Remote || opts.RemoteOnly) {
		err = git.DeleteRemoteBranch(branch.Name)
		if err != nil {
			logger.WarningF("Failed to delete remote branch \"%v\", %v\n", branch.Name, err)
		} else {
			logger.InfoF("Remote branch \"%v\" was deleted\n", branch.Name)
		}
	}
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolP("safe-delete", "s", false, "Safe delete branches - prevents deleting unmerged branches")
	deleteCmd.Flags().BoolP("ignore-errors", "i", false, "Ignore if git command fails and proceeds to remove the alias from the repository")
	deleteCmd.Flags().BoolP("remote", "r", false, "Delete the remote branch as well")
	deleteCmd.Flags().Bool("remote-only", false, "Deletes only the remote branch. Local branch and registry entry are not removed")
	deleteCmd.Flags().BoolP("worktree", "w", false, "Also delete the associated worktree")
	deleteCmd.Flags().BoolP("regex", "e", false, "Treat the arguments as regular expressions matched against branch names and aliases")
	deleteCmd.Flags().BoolP("all", "a", false, "Also consider local git branches not registered with g, expanding the candidate list before searching/deleting")
}
