/*
Copyright © 2024 Cyrus Mobini
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/spf13/cobra"
)

// getHomeCmd represents the getHome command
var getHomeCmd = &cobra.Command{
	Use:     "getHome",
	Short:   "Get the gitBranchTool's home directory path",
	Long:    `Get the gitBranchTool's home directory path`,
	Aliases: []string{"get-home", "home", "gh"},
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			internal.Logger.Fatal(err)
		}
		home = filepath.Join(home, internal.HOME_NAME)
		fmt.Println(home)
	},
}

func init() {
	rootCmd.AddCommand(getHomeCmd)
}
