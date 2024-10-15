package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const GITHUB_RELEASES = "https://github.com/cyrus2281/gitBranchTool/releases"

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
		Logger.CheckFatal(err)
		gHome = filepath.Join(home, HOME_NAME)
		if err := AddConfig("GIT_BRANCH_TOOL_HOME", gHome); err != nil {
			Logger.Fatal(err)
		}
	}
	return gHome
}

func GetRepositoryBranches() RepositoryBranches {
	git := Git{}
	// Get the home directory
	gHome := GetHome()
	repositoryName, err := git.GetRepositoryName()
	Logger.CheckFatal(err)
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

func FetchLatestVersion() (string, error) {
	url := "https://raw.githubusercontent.com/cyrus2281/gitBranchTool/refs/heads/main/VERSION"
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(body)), nil
}
