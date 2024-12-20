/*
Copyright © 2024 Cyrus Mobini
*/
package cmd

import (
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
		hasAlias := len(args) > 1
		var checkoutErr error

		repoBranches := internal.GetRepositoryBranches()
		branch, ok := repoBranches.GetBranchByNameOrAlias(id)
		if !ok {
			// branch doesn't exist
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

		if hasAlias && !ok {
			alias := args[1]
			notes := ""
			if len(args) > 2 {
				for _, note := range args[2:] {
					notes += note + " "
				}
			}

			newBranch := internal.Branch{
				Name:  id,
				Alias: alias,
				Note:  notes,
			}

			if repoBranches.AliasExists(newBranch.Alias) {
				logger.WarningF("Alias \"%v\" already exists. Alias must be unique\n", newBranch.Alias)
				return
			}

			repoBranches.AddBranch(newBranch)
			logger.InfoF("Branch \"%v\" has been registered with alias  \"%v\"\n", newBranch.Name, newBranch.Alias)
		} else if !hasAlias && !ok {
			logger.WarningF("Branch \"%v\" is not registered with alias\n", id)
		}
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
