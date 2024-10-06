/*
Copyright © 2024 Cyrus Mobini
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/cyrus2281/gitBranchTool/internal"
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
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		git := internal.Git{}
		if !git.IsGitRepo() {
			log.Fatalln("Not a git repository")
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
			fmt.Printf("Failed to switch branch to \"%v\"\n", id)
			log.Fatalln(checkoutErr)
		}
		fmt.Printf("Switched to branch \"%v\"\n", id)

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
				log.Printf("Alias \"%v\" already exists. Alias must be unique\n", newBranch.Alias)
				return
			}

			repoBranches.AddBranch(newBranch)
			fmt.Printf("Branch \"%v\" has been registered with alias  \"%v\"\n", newBranch.Name, newBranch.Alias)
		} else if !hasAlias && !ok {
			log.Printf("Branch \"%v\" is not registered with alias\n", id)
		}
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
	switchCmd.Aliases = []string{"checkout", "check", "s"}
}
