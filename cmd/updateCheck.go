/*
Copyright © 2024 Cyrus Mobini
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

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

func checkVersionAndUpdate(yesToAll bool) {
	latestVersion, downloadURL, err := checkLatestRelease()
	if err != nil {
		logger.Fatalln(err)
	}

	logger.InfoF("Current version: %s, Latest version: %s\n", rootCmd.Version, latestVersion)
	currentParts := strings.Split(rootCmd.Version, ".")
	latestParts := strings.Split(latestVersion, ".")

	if len(currentParts) != 3 || len(latestParts) != 3 {
		logger.Fatalln("Version format error")
	}

	if latestParts[0] > currentParts[0] {
		logger.Infoln("Major version update available. Please manually upgrade.")
		logger.InfoF("\tLatest release page: https://www.github.com/%s/releases/tag/V%s\n", internal.GITHUB_REPOSITORY, latestVersion)
		return
	} else if latestParts[1] > currentParts[1] || latestParts[2] > currentParts[2] {
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
	} else if latestVersion == rootCmd.Version {
		logger.Infoln("You're already on the latest version.")
	} else {
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
