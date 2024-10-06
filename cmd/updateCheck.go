/*
Copyright © 2024 Cyrus Mobini
*/
package cmd

import (
	"log"

	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/spf13/cobra"
)

// updateCheckCmd represents the updateCheck command
var updateCheckCmd = &cobra.Command{
	Use:     "updateCheck",
	Short:   "Checks if a newer version is available",
	Long:    `Checks if a newer version is available",`,
	Aliases: []string{"update-check", "uc"},
	Run: func(cmd *cobra.Command, args []string) {
		latestVersion, err := internal.FetchLatestVersion()
		if err != nil {
			log.Fatalf("Error fetching latest version: %v\n", err)
		}
		currentVersion := rootCmd.Version
		if latestVersion != currentVersion {
			log.Printf("You are running %v, a newer version is available: %v\n", currentVersion, latestVersion)
			log.Printf("Check the GitHub releases page for more information: %v\n", internal.GITHUB_RELEASES)
		} else {
			log.Printf("You are running the latest version: %v\n", currentVersion)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCheckCmd)
}
