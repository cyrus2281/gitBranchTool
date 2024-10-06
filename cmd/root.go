/*
Copyright © 2024 Cyrus Mobini
*/
package cmd

import (
	"os"
	"path/filepath"

	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "g",
	Short:   "A bash tool to facilitate managing git branches with long cryptic names with aliases",
	Long:    `A bash tool to facilitate managing git branches with long cryptic names with aliases`,
	Version: "3.0.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		home = filepath.Join(home, internal.HOME_NAME)
		if _, err := os.Stat(home); os.IsNotExist(err) {
			// Create the directory
			err = os.Mkdir(home, 0755)
			if err != nil {
				cobra.CheckErr(err)
			}
		}

		// Search config in home directory with name ".gitBranchTool" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("gitBranchTool.config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()
}
