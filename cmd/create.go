package cmd

import (
	"strings"

	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/cyrus2281/go-logger"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create NAME ALIAS [...NOTE]",
	Short: "Creates a branch with name, alias, and note, and checks into it",
	Long: `Creates a branch with name, alias, and note, and checks into it.
	Without only-create flag: "git checkout -b NAME"
	With only-create flag: "git branch NAME"
	`,
	Aliases: []string{"c"},
	Args:    cobra.MinimumNArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		alias := args[1]
		notes := args[2:]
		notesString := ""
		for _, note := range notes {
			notesString += note + " "
		}

		git := internal.Git{}
		if !git.IsGitRepo() {
			logger.Fatalln("Not a git repository")
		}

		repoBranches := internal.GetRepositoryBranches()

		prefix := repoBranches.GetLocalPrefix()
		if prefix == "" {
			prefix = internal.GetConfig(internal.GLOBAL_PREFIX_KEY)
		}

		worktreeAlias, _ := cmd.Flags().GetString("worktree")
		useWorktree := cmd.Flags().Changed("worktree")
		worktreeAlias = strings.TrimSpace(worktreeAlias)
		createOnly, _ := cmd.Flags().GetBool("only-create")

		opts := createOpts{
			Prefix:   prefix,
			CreateOnly: createOnly,
			UseWorktree: useWorktree,
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

		executeCreate(&git, &repoBranches, name, alias, notesString, opts)
	},
}

type createOpts struct {
	Prefix               string
	CreateOnly           bool
	UseWorktree          bool
	WorktreeAlias        string
	WorktreePathTemplate string
	RepoPath             string
	RepoName             string
}

func executeCreate(git *internal.Git, repoBranches *internal.RepositoryBranches,
	name, alias, notes string, opts createOpts) {

	if opts.Prefix != "" {
		name = opts.Prefix + name
	}

	newBranch := internal.Branch{
		Name:  name,
		Alias: alias,
		Note:  notes,
	}
	if repoBranches.AliasExists(newBranch.Alias) {
		logger.Fatalln("Alias already exists. Alias must be unique.")
	}
	if repoBranches.BranchWithAliasExists(newBranch.Alias) {
		logger.FatalF("A branch with name \"%s\" already exists. Alias must be unique.\n", newBranch.Alias)
	}

	if opts.UseWorktree {
		worktreeAlias := opts.WorktreeAlias
		if worktreeAlias == "" {
			worktreeAlias = alias
		}

		if repoBranches.WorktreeAliasExists(worktreeAlias) {
			logger.FatalF("Worktree alias \"%s\" already exists. Alias must be unique.\n", worktreeAlias)
		}

		resolvedPath := internal.ResolveWorktreePath(opts.WorktreePathTemplate, opts.RepoPath, opts.RepoName, worktreeAlias, name)

		err := git.WorktreeAddNewBranch(name, resolvedPath)
		if err != nil {
			logger.Fatalln("Failed to create worktree:", err)
		}

		repoBranches.AddBranch(newBranch)
		logger.InfoF("Branch %v with alias %v was created\n", newBranch.Name, newBranch.Alias)

		newWorktree := internal.Worktree{
			Alias: worktreeAlias,
			Path:  resolvedPath,
			Note:  notes,
		}
		repoBranches.AddWorktree(newWorktree)
		logger.InfoF("Worktree \"%s\" created at %s\n", worktreeAlias, resolvedPath)
	} else {
		var err error
		if opts.CreateOnly {
			err = git.CreateNewBranch(newBranch.Name)
		} else {
			err = git.SwitchToNewBranch(newBranch.Name)
		}
		if err != nil {
			logger.Fatalln(err)
		}

		repoBranches.AddBranch(newBranch)
		logger.InfoF("Branch %v with alias %v was created\n", newBranch.Name, newBranch.Alias)
	}
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().BoolP("only-create", "o", false, "Only create the branch, do not check into it")
	createCmd.Flags().StringP("worktree", "w", "", "Also create a worktree for the branch. Use -w=ALIAS to specify a custom worktree alias.")
	createCmd.Flags().Lookup("worktree").NoOptDefVal = " "
}
