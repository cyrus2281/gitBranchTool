package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cyrus2281/go-logger"
	"github.com/spf13/viper"
)

const GITHUB_REPOSITORY = "cyrus2281/gitBranchTool"

const HOME_NAME = ".gitBranchTool_go"

const CONFIG_NAME = "gitBranchTool.config"

func AddConfig(key string, value any) error {
	viper.Set(key, value)
	err := viper.SafeWriteConfig()
	return err
}

func GetHome() string {
	gHome := viper.GetString("GIT_BRANCH_TOOL_HOME")
	if gHome == "" {
		home, err := os.UserHomeDir()
		logger.CheckFatalln(err)
		gHome = filepath.Join(home, HOME_NAME)
		if err := AddConfig("GIT_BRANCH_TOOL_HOME", gHome); err != nil {
			logger.Fatalln(err)
		}
	}
	return gHome
}

func GetRepositoryBranches() RepositoryBranches {
	git := Git{}
	// Get the home directory
	gHome := GetHome()
	repositoryName, err := git.GetRepositoryName()
	logger.CheckFatalln(err)
	repoBranches := RepositoryBranches{
		RepositoryName: repositoryName,
		StoreDirectory: gHome,
	}
	return repoBranches
}

func PrintTableHeader() {
	fmt.Printf("   %-20s\t%-20s\t%-20s\n", "Branch Name", "Alias", "Note")
	fmt.Println("------------------------------------------------------------")
}
