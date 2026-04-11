package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/cyrus2281/go-logger"
)

func saveUpdateCheckTimestamp() {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	if err := internal.AddConfig(internal.LAST_UPDATE_CHECK_KEY, timestamp); err != nil {
		logger.Debugln("Failed to save update check timestamp:", err)
	}
}

func checkUpdateReminder() {
	lastCheck := internal.GetConfig(internal.LAST_UPDATE_CHECK_KEY)
	if lastCheck != "" {
		ts, err := strconv.ParseInt(lastCheck, 10, 64)
		if err == nil {
			elapsed := time.Since(time.Unix(ts, 0))
			if elapsed < time.Duration(internal.UPDATE_CHECK_INTERVAL_DAYS)*24*time.Hour {
				return
			}
		}
	}

	latestVersion, err := checkLatestReleaseWithTimeout(2)
	if err != nil {
		logger.Debugln("Update reminder check failed:", err)
		return
	}

	saveUpdateCheckTimestamp()

	status, err := compareVersions(rootCmd.Version, latestVersion)
	if err != nil {
		logger.Debugln("Version comparison failed:", err)
		return
	}

	if status == VersionMinorUpdate || status == VersionMajorUpdate {
		fmt.Printf("Update available: v%s (current: v%s). Run 'g updateCheck' to upgrade.\n", latestVersion, rootCmd.Version)
	}
}
