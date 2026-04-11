package internal

import (
	"os"
	"path/filepath"

	"github.com/cyrus2281/go-logger"
	"github.com/spf13/viper"
)

const GITHUB_REPOSITORY = "cyrus2281/gitBranchTool"

const HOME_NAME = ".gitBranchTool_go"

const CONFIG_NAME = "gitBranchTool.config"

const WORKTREE_PATH_KEY = "GIT_BRANCH_TOOL_WORKTREE_PATH"
const DELETE_BRANCHES_WORKTREE_KEY = "GIT_BRANCH_TOOL_DELETE_BRANCHES_WORKTREE"
const DEFAULT_WORKTREE_PATH = "../{repository}-worktrees/{alias}"

func AddConfig(key string, value any) error {
	viper.Set(key, value)
	err := viper.WriteConfig()
	return err
}

func GetConfig(key string) string {
	return viper.GetString(key)
}

func GetHome() string {
	gHome := viper.GetString("GIT_BRANCH_TOOL_HOME")
	if gHome == "" {
		home, err := os.UserHomeDir()
		logger.CheckFatalln(err)
		gHome = filepath.Join(home, HOME_NAME)
		viper.SafeWriteConfig()
		if err := AddConfig("GIT_BRANCH_TOOL_HOME", gHome); err != nil {
			logger.Fatalln(err)
		}
	}
	logger.Debugln("Home: ", gHome)
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
	logger.InfoF("   %-20s\t%-20s\t%-20s\n", "Branch Name", "Alias", "Note")
	logger.Infoln("------------------------------------------------------------")
}

func PrintBranchTableHeaderWithWorktree() {
	logger.InfoF("   %-20s\t%-20s\t%-20s\t%-15s\n", "Branch Name", "Alias", "Note", "Worktree")
	logger.Infoln("-------------------------------------------------------------------------------")
}

func PrintWorktreeTableHeader() {
	logger.InfoF("   %-40s\t%-15s\t%-25s\t%-15s\t%-20s\n", "Path", "Alias", "Branch", "Branch Alias", "Note")
	logger.Infoln("------------------------------------------------------------------------------------------------------------------------------")
}

func GetWorktreePath() string {
	path := GetConfig(WORKTREE_PATH_KEY)
	if path == "" {
		return DEFAULT_WORKTREE_PATH
	}
	return path
}
