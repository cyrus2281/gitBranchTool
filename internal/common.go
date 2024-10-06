package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

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
		if err != nil {
			log.Fatalln(err)
		}
		gHome = filepath.Join(home, HOME_NAME)
		if err := AddConfig("GIT_BRANCH_TOOL_HOME", gHome); err != nil {
			log.Fatalln(err)
		}
	}
	return gHome
}

func GetRepositoryBranches() RepositoryBranches {
	git := Git{}
	// Get the home directory
	gHome := GetHome()
	repositoryName, err := git.GetRepositoryName()
	if err != nil {
		log.Fatalln(err)
	}
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
