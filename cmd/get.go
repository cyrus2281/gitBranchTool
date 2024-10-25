/*
Copyright © 2024 Cyrus Mobini
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/cyrus2281/go-logger"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get COMMAND",
	Short: "Get configuration options (Run `g get --help` for more information)",
	Long: `Get configuration options

- get default-branch: Get the default branch for the current repository

- get local-prefix: Get the branch prefix for the current repository

- get global-prefix: Get the branch prefix for all repositories

- get home: Get the path to the home directory for the g command
`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		position := len(args) + 1
		commands := []string{"default-branch", "local-prefix", "global-prefix", "home"}
		if position == 1 {
			return commands, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		command := args[0]
		switch command {
		case "default-branch":
			git := internal.Git{}
			if !git.IsGitRepo() {
				logger.Fatalln("Not a git repository")
			}
			repoBranches := internal.GetRepositoryBranches()
			name := repoBranches.GetDefaultBranch()
			logger.InfoF("Default branch for the repository \"%v\" is ", repoBranches.RepositoryName)
			fmt.Println(name)
		case "local-prefix":
			git := internal.Git{}
			if !git.IsGitRepo() {
				logger.Fatalln("Not a git repository")
			}
			repoBranches := internal.GetRepositoryBranches()
			prefix := repoBranches.GetLocalPrefix()
			if prefix == "" {
				logger.InfoF("No local prefix set for the repository \"%v\"\n", repoBranches.RepositoryName)
				return
			} else {
				logger.InfoF("Local prefix for the repository \"%v\" is", repoBranches.RepositoryName)
			}
			fmt.Println(prefix)
		case "global-prefix":
			prefix := internal.GetConfig(internal.GLOBAL_PREFIX_KEY)
			if prefix == "" {
				logger.InfoF("No global prefix set\n")
			} else {
				logger.InfoF("Global branch prefix is ")
			}
			fmt.Println(prefix)
		case "home":
			home, err := os.UserHomeDir()
			if err != nil {
				logger.Fatalln(err)
			}
			home = filepath.Join(home, internal.HOME_NAME)
			logger.InfoF("Home directory is ")
			fmt.Println(home)
		default:
			logger.Fatalln("Invalid command!\n\tAvailable commands: default-branch, local-prefix, global-prefix, home")
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
