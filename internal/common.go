package internal

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func GetHome() string {
	gHome := viper.GetString("GIT_BRANCH_TOOL_HOME")
	if gHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		gHome = filepath.Join(home, ".gitBranchTool_go")
		viper.Set("GIT_BRANCH_TOOL_HOME", gHome)
		if err := viper.SafeWriteConfig(); err != nil {
			panic(err)
		}
	}
	return gHome
}
