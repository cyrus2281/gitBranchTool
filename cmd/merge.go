package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/cyrus2281/go-logger"
	"github.com/spf13/cobra"
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge [NAME/ALIAS]",
	Short: "Merges or rebases the given branch into the current branch",
	Long: `Merges (or rebases) the branch with the given name or alias into the current branch.
Uses the git command "git merge NAME" (or "git rebase NAME" with --rebase).

If no name or alias is given, the configured default branch is used. This makes
"g merge" a quick way to update the current branch with the latest of main.

When a merge or rebase stops due to conflicts, resolve them and run
"g m --continue" to finish, or "g m --abort" to cancel.

Examples:
	g merge main            Merge the default/main branch into the current branch
	g merge feat -r         Rebase the current branch onto the "feat" branch
	g merge main -f         Fetch origin/main first, then merge it
	g merge --continue      Continue after resolving conflicts`,
	Args:    cobra.MaximumNArgs(1),
	Aliases: []string{"m"},
	Annotations: map[string]string{
		manualAnnotation: `Merge (or rebase) the branch with the given NAME/ALIAS into the current branch. With no argument, uses the configured default branch (a quick way to sync the current branch with main).
Flags: -r/--rebase, -s/--squash, -f/--fetch (fetch from origin first), --ff-only, --no-ff, -n/--no-verify, --continue and --abort (for an in-progress merge/rebase).`,
	},
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

		opts := mergeOpts{Remote: "origin"}
		opts.Rebase, _ = cmd.Flags().GetBool("rebase")
		opts.Squash, _ = cmd.Flags().GetBool("squash")
		opts.NoVerify, _ = cmd.Flags().GetBool("no-verify")
		opts.Fetch, _ = cmd.Flags().GetBool("fetch")
		opts.FFOnly, _ = cmd.Flags().GetBool("ff-only")
		opts.NoFF, _ = cmd.Flags().GetBool("no-ff")
		opts.Abort, _ = cmd.Flags().GetBool("abort")
		opts.Continue, _ = cmd.Flags().GetBool("continue")

		if err := validateMergeOpts(opts, len(args) > 0); err != nil {
			logger.Fatalln(err)
		}

		if opts.Abort || opts.Continue {
			executeMergeState(&git, opts)
			return
		}

		repoBranches := internal.GetRepositoryBranches()

		id := ""
		if len(args) > 0 {
			id = args[0]
		} else {
			id = repoBranches.GetDefaultBranch()
			logger.InfoF("No branch given; using default branch \"%s\"\n", id)
		}

		executeMerge(&git, &repoBranches, id, opts)
	},
}

type mergeOpts struct {
	Rebase   bool
	Squash   bool
	NoVerify bool
	Fetch    bool
	FFOnly   bool
	NoFF     bool
	Abort    bool
	Continue bool
	Remote   string
}

// validateMergeOpts enforces the flag-combination rules. It is a pure function
// (no I/O) so it can be unit-tested directly.
func validateMergeOpts(opts mergeOpts, hasArg bool) error {
	if opts.Abort && opts.Continue {
		return fmt.Errorf("cannot combine --abort and --continue")
	}
	if (opts.Abort || opts.Continue) && hasArg {
		return fmt.Errorf("--abort/--continue do not take a branch argument")
	}
	if opts.Rebase {
		var bad []string
		if opts.Squash {
			bad = append(bad, "--squash")
		}
		if opts.NoVerify {
			bad = append(bad, "--no-verify")
		}
		if opts.FFOnly {
			bad = append(bad, "--ff-only")
		}
		if opts.NoFF {
			bad = append(bad, "--no-ff")
		}
		if len(bad) > 0 {
			return fmt.Errorf("--rebase cannot be combined with %s", strings.Join(bad, ", "))
		}
	}
	if opts.FFOnly && opts.NoFF {
		return fmt.Errorf("cannot combine --ff-only and --no-ff")
	}
	return nil
}

// buildMergeArgs assembles the extra flags passed to "git merge".
func buildMergeArgs(opts mergeOpts) []string {
	var args []string
	if opts.Squash {
		args = append(args, "--squash")
	}
	if opts.NoVerify {
		args = append(args, "--no-verify")
	}
	if opts.FFOnly {
		args = append(args, "--ff-only")
	}
	if opts.NoFF {
		args = append(args, "--no-ff")
	}
	return args
}

// executeMerge resolves the target branch and performs the merge or rebase.
func executeMerge(git *internal.Git, repoBranches *internal.RepositoryBranches, id string, opts mergeOpts) {
	branch, registered := repoBranches.GetBranchByNameOrAlias(id)
	branchName := id
	if registered {
		branchName = branch.Name
	}

	// Same-branch guard. Skipped when fetching, since fetching the current
	// branch's remote and merging it is a valid "pull" operation.
	if !opts.Fetch {
		if current, err := git.GetCurrentBranch(); err == nil && current == branchName {
			logger.InfoF("Already on branch \"%s\"; nothing to merge\n", branchName)
			return
		}
	}

	ref := branchName
	if opts.Fetch {
		remote := opts.Remote
		if remote == "" {
			remote = "origin"
		}
		logger.InfoF("Fetching \"%s\" from \"%s\"...\n", branchName, remote)
		if err := git.Fetch(remote, branchName); err != nil {
			logger.ErrorF("Failed to fetch \"%s\" from \"%s\"\n", branchName, remote)
			logger.Fatalln(err)
		}
		ref = "FETCH_HEAD"
	}

	if opts.Rebase {
		if err := git.RebaseBranch(ref); err != nil {
			handleIntegrationFailure(git, "rebase", branchName)
		}
		logger.InfoF("Rebased the current branch onto \"%s\"\n", branchName)
		return
	}

	if err := git.MergeBranch(ref, buildMergeArgs(opts)); err != nil {
		handleIntegrationFailure(git, "merge", branchName)
	}
	if opts.Squash {
		logger.InfoF("Squash-merged \"%s\". Changes are staged; commit them to finish.\n", branchName)
		return
	}
	logger.InfoF("Merged \"%s\" into the current branch\n", branchName)
}

// handleIntegrationFailure reports a failed merge/rebase and exits non-zero.
// When the operation is left in progress it is treated as a conflict and the
// user is told how to continue or abort.
func handleIntegrationFailure(git *internal.Git, action, branchName string) {
	if git.MergeInProgress() || git.RebaseInProgress() {
		logger.WarningF("Conflicts while trying to %s \"%s\".\n", action, branchName)
		logger.WarningF("Resolve them, then run \"g m --continue\" (or \"g m --abort\" to cancel).\n")
		os.Exit(1)
	}
	logger.ErrorF("Failed to %s \"%s\"\n", action, branchName)
	os.Exit(1)
}

// executeMergeState handles --abort and --continue for an in-progress
// merge or rebase, auto-detecting which one is running.
func executeMergeState(git *internal.Git, opts mergeOpts) {
	rebase := git.RebaseInProgress()
	merge := git.MergeInProgress()
	if !rebase && !merge {
		logger.Fatalln("No merge or rebase in progress")
	}

	word := "merge"
	if rebase {
		word = "rebase"
	}

	if opts.Abort {
		var err error
		if rebase {
			err = git.RebaseAbort()
		} else {
			err = git.MergeAbort()
		}
		if err != nil {
			logger.Fatalln(err)
		}
		logger.InfoF("Aborted the in-progress %s\n", word)
		return
	}

	// --continue
	var err error
	if rebase {
		err = git.RebaseContinue()
	} else {
		err = git.MergeContinue()
	}
	if err != nil {
		if git.MergeInProgress() || git.RebaseInProgress() {
			logger.WarningF("There are still unresolved conflicts. Resolve them and run \"g m --continue\" again.\n")
			os.Exit(1)
		}
		logger.Fatalln(err)
	}
	logger.InfoF("Continued and completed the %s\n", word)
}

func init() {
	rootCmd.AddCommand(mergeCmd)
	mergeCmd.Flags().BoolP("rebase", "r", false, "Use git rebase instead of git merge")
	mergeCmd.Flags().BoolP("squash", "s", false, "Squash the merged commits (git merge --squash; incompatible with --rebase)")
	mergeCmd.Flags().BoolP("no-verify", "n", false, "Skip git hooks (git merge --no-verify; incompatible with --rebase)")
	mergeCmd.Flags().BoolP("fetch", "f", false, "Fetch the latest of the branch from origin before merging/rebasing")
	mergeCmd.Flags().Bool("ff-only", false, "Refuse to merge unless it can be fast-forwarded (incompatible with --rebase)")
	mergeCmd.Flags().Bool("no-ff", false, "Always create a merge commit (incompatible with --rebase)")
	mergeCmd.Flags().Bool("abort", false, "Abort an in-progress merge or rebase")
	mergeCmd.Flags().Bool("continue", false, "Continue an in-progress merge or rebase after resolving conflicts")
}
