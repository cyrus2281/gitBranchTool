/*
Copyright © 2024 Cyrus Mobini
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// getHomeCmd represents the getHome command
var getHomeCmd = &cobra.Command{
	Use:   "getHome",
	Short: "Get the gitBranchTool's home directory path",
	Long:  `Get the gitBranchTool's home directory path`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		home = filepath.Join(home, ".gitBranchTool_go")
		fmt.Println(home)
	},
}

func init() {
	rootCmd.AddCommand(getHomeCmd)
	getHomeCmd.Aliases = []string{"home", "gh"}

}
