package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/cyrus2281/gitBranchTool/internal"
	"github.com/cyrus2281/go-logger"
	"github.com/spf13/cobra"
)

// updateCheckCmd represents the updateCheck command
var updateCheckCmd = &cobra.Command{
	Use:     "updateCheck",
	Short:   "Checks if a newer version is available",
	Long:    `Checks if a newer version is available. Asks to upgrade if available",`,
	Aliases: []string{"update-check", "uc"},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		yesToAll, _ := cmd.Flags().GetBool("yes-to-all")
		checkVersionAndUpdate(yesToAll)
		saveUpdateCheckTimestamp()
	},
}

func init() {
	rootCmd.AddCommand(updateCheckCmd)
	updateCheckCmd.Flags().BoolP("yes-to-all", "y", false, "Automatically update to the latest version")
}

type Release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

type VersionStatus int

const (
	VersionUnofficial  VersionStatus = -1
	VersionUpToDate    VersionStatus = 0
	VersionMinorUpdate VersionStatus = 1
	VersionMajorUpdate VersionStatus = 2
)

func compareVersions(current, latest string) (VersionStatus, error) {
	currentParts := strings.Split(current, ".")
	latestParts := strings.Split(latest, ".")

	if len(currentParts) != 3 || len(latestParts) != 3 {
		return 0, fmt.Errorf("version format error")
	}

	currentNums := make([]int, 3)
	latestNums := make([]int, 3)
	for i := 0; i < 3; i++ {
		var err error
		currentNums[i], err = strconv.Atoi(currentParts[i])
		if err != nil {
			return 0, fmt.Errorf("version format error: %w", err)
		}
		latestNums[i], err = strconv.Atoi(latestParts[i])
		if err != nil {
			return 0, fmt.Errorf("version format error: %w", err)
		}
	}

	if latestNums[0] > currentNums[0] {
		return VersionMajorUpdate, nil
	}
	if latestNums[0] < currentNums[0] {
		return VersionUnofficial, nil
	}
	// Major equal
	if latestNums[1] > currentNums[1] {
		return VersionMinorUpdate, nil
	}
	if latestNums[1] < currentNums[1] {
		return VersionUnofficial, nil
	}
	// Minor equal
	if latestNums[2] > currentNums[2] {
		return VersionMinorUpdate, nil
	}
	if latestNums[2] < currentNums[2] {
		return VersionUnofficial, nil
	}
	return VersionUpToDate, nil
}

func checkVersionAndUpdate(yesToAll bool) {
	latestVersion, downloadURL, err := checkLatestRelease()
	if err != nil {
		logger.Fatalln(err)
	}

	logger.InfoF("Current version: %s, Latest version: %s\n", rootCmd.Version, latestVersion)

	status, err := compareVersions(rootCmd.Version, latestVersion)
	if err != nil {
		logger.Fatalln(err)
	}

	switch status {
	case VersionMajorUpdate:
		logger.Infoln("Major version update available. Please manually upgrade.")
		logger.InfoF("\tLatest release page: https://www.github.com/%s/releases/tag/V%s\n", internal.GITHUB_REPOSITORY, latestVersion)
	case VersionMinorUpdate:
		downloadLatest := false
		if yesToAll {
			logger.Infoln("Updating to the latest version...")
			downloadLatest = true
		} else {
			logger.InfoF("New version available: %s. Would you like to update? (y/[n]): ", latestVersion)
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.ToLower(strings.TrimSpace(input))
			if input == "y" || input == "yes" {
				downloadLatest = true
			}
		}
		if downloadLatest {
			if err := downloadAndReplace(downloadURL, latestVersion); err != nil {
				logger.Fatalln("Error updating: ", err)
			}
		}
	case VersionUpToDate:
		logger.Infoln("You're already on the latest version.")
	case VersionUnofficial:
		logger.Infoln("You're on an unofficial version. Please check the latest release.")
		logger.InfoF("\tLatest release page: https://www.github.com/%s/releases/tag/V%s\n", internal.GITHUB_REPOSITORY, latestVersion)
	}
}

func checkLatestRelease() (string, string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", internal.GITHUB_REPOSITORY)
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("failed to fetch release info: status code %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", "", fmt.Errorf("failed to decode response: %w", err)
	}

	latestVersion := strings.TrimPrefix(strings.ToLower(release.TagName), "v")

	downloadURL := ""
	for _, asset := range release.Assets {
		if strings.HasPrefix(asset.Name, getAssetPrefix(latestVersion)) {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return "", "", fmt.Errorf("no suitable version found for your OS")
	}

	return latestVersion, downloadURL, nil
}

func checkLatestReleaseWithTimeout(seconds int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	defer cancel()

	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", internal.GITHUB_REPOSITORY)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch release info: status code %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	latestVersion := strings.TrimPrefix(strings.ToLower(release.TagName), "v")
	return latestVersion, nil
}

func getAssetPrefix(version string) string {
	version = strings.TrimPrefix(version, "v")
	switch runtime.GOOS {
	case "linux":
		return fmt.Sprintf("g-linux-v%s", version)
	case "darwin":
		return fmt.Sprintf("g-macos-v%s", version)
	case "windows":
		return fmt.Sprintf("g-win-v%s.exe", version)
	default:
		return ""
	}
}

func downloadAndReplace(downloadURL, version string) error {
	resp, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to download version %s: status code %d", version, resp.StatusCode)
	}

	tempFile, err := os.CreateTemp("", fmt.Sprintf("g-v%s-*", version))
	if err != nil {
		return err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return err
	}

	// Find current executable path
	executablePath, err := os.Executable()
	if err != nil {
		return err
	}

	// If the OS is not Windows, update the binary directly
	if runtime.GOOS != "windows" {
		return nonWindowBinaryUpdate(tempFile.Name(), executablePath)
	} else {
		// PowerShell doesn't allow renaming the running executable
		return windowBinaryUpdate(tempFile.Name(), executablePath)
	}
}

func nonWindowBinaryUpdate(tempPath, executablePath string) error {
	if err := os.Rename(tempPath, executablePath); err != nil {
		// If the error has permission denied
		if strings.Contains(err.Error(), "permission denied") {
			return fmt.Errorf("permission denied. Please run with sudo")
		}
		return err
	}

	// Set executable permission for Unix-based systems
	if err := os.Chmod(executablePath, 0755); err != nil {
		return err
	}

	logger.Infoln("Successfully updated to latest version")
	return nil
}

func windowBinaryUpdate(tempPath, executablePath string) error {
	// PowerShell doesn't allow renaming the running executable
	// So we need to create a batch script to rename the executable
	batchScript := fmt.Sprintf(`@echo off
timeout /t 2 /nobreak > NUL
move /Y "%s" "%s"
exit
`, tempPath, executablePath)

	batchPath := executablePath + ".bat"
	if err := os.WriteFile(batchPath, []byte(batchScript), 0755); err != nil {
		return err
	}

	// Run the batch script
	cmd := exec.Command("cmd", "/C", batchPath)
	if err := cmd.Start(); err != nil {
		return err
	}

	logger.Infoln("Successfully updated to latest version")
	return nil
}
