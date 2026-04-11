package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cyrus2281/go-logger"
	"github.com/spf13/viper"
)

const GITHUB_REPOSITORY = "cyrus2281/gitBranchTool"

const HOME_NAME = ".gitBranchTool_go"

const CONFIG_NAME = "gitBranchTool.config"

const WORKTREE_PATH_KEY = "GIT_BRANCH_TOOL_WORKTREE_PATH"
const DELETE_BRANCHES_WORKTREE_KEY = "GIT_BRANCH_TOOL_DELETE_BRANCHES_WORKTREE"
const DEFAULT_WORKTREE_PATH = "./worktrees/{alias}"

const LAST_UPDATE_CHECK_KEY = "GIT_BRANCH_TOOL_LAST_UPDATE_CHECK"
const UPDATE_CHECK_INTERVAL_DAYS = 5

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

// PrintTable prints a formatted table with dynamic column widths.
// Headers define column names, rows contain the data.
// Row numbering (0), 1), ...) is added automatically.
func PrintTable(headers []string, rows [][]string) {
	colCount := len(headers)

	// Compute max width per column from headers and data
	widths := make([]int, colCount)
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i := 0; i < colCount && i < len(row); i++ {
			if len(row[i]) > widths[i] {
				widths[i] = len(row[i])
			}
		}
	}

	// Compute index prefix width (e.g. "0) " = 3, "10) " = 4)
	indexWidth := 3
	if len(rows) >= 10 {
		indexWidth = len(fmt.Sprintf("%d) ", len(rows)-1))
	}
	padding := 3

	// Print header
	fmt.Printf("%s", strings.Repeat(" ", indexWidth))
	for i, h := range headers {
		if i < colCount-1 {
			fmt.Printf("%-*s%s", widths[i], h, strings.Repeat(" ", padding))
		} else {
			fmt.Printf("%s", h)
		}
	}
	fmt.Println()

	// Print separator
	totalWidth := indexWidth
	for i, w := range widths {
		totalWidth += w
		if i < colCount-1 {
			totalWidth += padding
		}
	}
	fmt.Println(strings.Repeat("-", totalWidth))

	// Print rows
	for idx, row := range rows {
		prefix := fmt.Sprintf("%d) ", idx)
		fmt.Printf("%-*s", indexWidth, prefix)
		for i := 0; i < colCount; i++ {
			cell := ""
			if i < len(row) {
				cell = row[i]
			}
			if i < colCount-1 {
				fmt.Printf("%-*s%s", widths[i], cell, strings.Repeat(" ", padding))
			} else {
				fmt.Printf("%s", cell)
			}
		}
		fmt.Println()
	}
}

func GetWorktreePath() string {
	path := GetConfig(WORKTREE_PATH_KEY)
	if path == "" {
		return DEFAULT_WORKTREE_PATH
	}
	return path
}
