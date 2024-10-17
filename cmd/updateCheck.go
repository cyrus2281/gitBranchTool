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
	"github.com/spf13/cobra"
)

// updateCheckCmd represents the updateCheck command
var updateCheckCmd = &cobra.Command{
	Use:     "updateCheck",
	Short:   "Checks if a newer version is available",
	Long:    `Checks if a newer version is available. Asks to upgrade if available",`,
	Aliases: []string{"update-check", "uc"},
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
		internal.Logger.Fatal(err)
	}

	internal.Logger.InfoF("Current version: %s, Latest version: %s\n", rootCmd.Version, latestVersion)
	currentParts := strings.Split(rootCmd.Version, ".")
	latestParts := strings.Split(latestVersion, ".")

	if len(currentParts) != 3 || len(latestParts) != 3 {
		internal.Logger.Fatal("Version format error")
	}

	if latestParts[0] > currentParts[0] {
		internal.Logger.Info("Major version update available. Please manually upgrade.")
		internal.Logger.InfoF("\tLatest release page: https://www.github.com/%s/releases/tag/V%s\n", internal.GITHUB_REPOSITORY, latestVersion)
		return
	} else if latestParts[1] > currentParts[1] || latestParts[2] > currentParts[2] {
		downloadLatest := false
		if yesToAll {
			internal.Logger.Info("Updating to the latest version...")
			downloadLatest = true
		} else {
			internal.Logger.InfoF("New version available: %s. Would you like to update? (y/[n]): ", latestVersion)
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.ToLower(strings.TrimSpace(input))
			if input == "y" || input == "yes" {
				downloadLatest = true
			}
		}

		if downloadLatest {
			if err := downloadAndReplace(downloadURL, latestVersion); err != nil {
				internal.Logger.Fatal("Error updating: ", err)
			} else {
				internal.Logger.InfoF("Successfully updated to version %s\n", latestVersion)
			}
		}
	} else if latestVersion == rootCmd.Version {
		internal.Logger.Info("You're already on the latest version.")
	} else {
		internal.Logger.Info("You're on an unofficial version. Please check the latest release.")
		internal.Logger.InfoF("\tLatest release page: https://www.github.com/%s/releases/tag/V%s\n", internal.GITHUB_REPOSITORY, latestVersion)
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
		return fmt.Sprintf("g.win-v%s.exe", version)
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
	whichG := exec.Command("which", "g")
	pathBytes, err := whichG.Output()
	if err != nil {
		return err
	}

	executablePath := strings.TrimSpace(string(pathBytes))

	// Rename and replace
	if runtime.GOOS == "windows" {
		executablePath += ".exe"
	}

	if err := os.Rename(tempFile.Name(), executablePath); err != nil {
		// If the error has permission denied
		if strings.Contains(err.Error(), "permission denied") {
			return fmt.Errorf("permission denied. Please run with sudo or as an administrator")
		}
		return err
	}

	// Set executable permission for Unix-based systems
	if runtime.GOOS != "windows" {
		if err := os.Chmod(executablePath, 0755); err != nil {
			return err
		}
	}

	return nil
}
