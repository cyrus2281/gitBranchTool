package cmd

import (
	"strconv"
	"strings"
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

func clearRequiresUpdate() {
	if err := internal.AddConfig(internal.REQUIRES_UPDATE_KEY, "false"); err != nil {
		logger.Debugln("Failed to clear requires_update:", err)
	}
}

func isCheckIntervalElapsed() bool {
	lastCheck := internal.GetConfig(internal.LAST_UPDATE_CHECK_KEY)
	if lastCheck == "" {
		return true
	}
	ts, err := strconv.ParseInt(lastCheck, 10, 64)
	if err != nil {
		return true
	}
	elapsed := time.Since(time.Unix(ts, 0))
	return elapsed >= time.Duration(internal.UPDATE_CHECK_INTERVAL_DAYS)*24*time.Hour
}

func checkUpdateReminder() {
	requiresUpdate := strings.ToLower(internal.GetConfig(internal.REQUIRES_UPDATE_KEY))
	if requiresUpdate == "true" {
		logger.InfoF("\n\tUpdate available. Run 'g updateCheck' to upgrade.\n")
	}

	if !isCheckIntervalElapsed() {
		return
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
		if err := internal.AddConfig(internal.REQUIRES_UPDATE_KEY, "true"); err != nil {
			logger.Debugln("Failed to set requires_update:", err)
		}
		if requiresUpdate != "true" {
			logger.InfoF("\n\tUpdate available: v%s (current: v%s). Run 'g updateCheck' to upgrade.\n", latestVersion, rootCmd.Version)
		}
	} else {
		clearRequiresUpdate()
	}
}
