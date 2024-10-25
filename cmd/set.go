/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/cyrus2281/go-logger"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set COMMAND [ARG]",
	Short: "Set configuration options (Run `g set --help` for more information)",
	Long: `Set configuration options

- set default-branch <NAME> : Change the default branch, default is "main"
	Example: g set default-branch master

- set local-prefix <PREFIX> : Change the branch prefix for the current repository, All branches created with the g command will have this prefix.
Default is nothing (Run command with no argument to remove prefix)(Overrides global-prefix)
	Example: g set local-prefix feature/

- set global-prefix <PREFIX> : Change the branch prefix for all repositories, All branches created with the g command will have this prefix.
Default is nothing (Run command with no argument to remove prefix)
	Example: g set global-prefix feature/
`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		position := len(args) + 1
		commands := []string{"default-branch", "local-prefix", "global-prefix"}
		if position == 1 {
			return commands, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		command := args[0]
		switch command {
		case "default-branch":
			git := internal.Git{}
			if !git.IsGitRepo() {
				logger.Fatalln("Not a git repository")
			}
			if len(args) != 2 {
				logger.Fatalln("Invalid number of arguments. Run `g set default-branch <NAME>`")
			}
			name := args[1]

			repoBranches := internal.GetRepositoryBranches()
			repoBranches.SetDefaultBranch(name)
			logger.InfoF("Default branch for the repository \"%v\" set to \"%v\"\n", repoBranches.RepositoryName, name)
		case "local-prefix":
			git := internal.Git{}
			if !git.IsGitRepo() {
				logger.Fatalln("Not a git repository")
			}

			prefix := ""
			if len(args) == 2 {
				prefix = args[1]
			} else if len(args) > 2 {
				logger.Fatalln("Invalid number of arguments. Run `g set local-prefix <PREFIX>`")
			}
			repoBranches := internal.GetRepositoryBranches()
			repoBranches.SetLocalPrefix(prefix)
			logger.InfoF("Local prefix for the repository \"%v\" set to \"%v\"\n", repoBranches.RepositoryName, prefix)
		case "global-prefix":
			prefix := ""
			if len(args) == 2 {
				prefix = args[1]
			} else if len(args) > 2 {
				logger.Fatalln("Invalid number of arguments. Run `g set local-prefix <PREFIX>`")
			}
			err := internal.AddConfig(internal.GLOBAL_PREFIX_KEY, prefix)
			logger.CheckFatalln(err)
			logger.InfoF("Global branch prefix set to \"%v\"\n", prefix)
		default:
			logger.Fatalln("Invalid command!\n\tAvailable commands: default-branch, local-prefix, global-prefix")
		}
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
