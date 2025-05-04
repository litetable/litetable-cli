package cmd

import (
	"bufio"
	"fmt"
	"github.com/litetable/litetable-cli/internal/litetable"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update LiteTable CLI",
	Long:  "Check for and install the latest version of LiteTable CLI",
	Run: func(cmd *cobra.Command, args []string) {
		if err := updateCLI(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(UpdateCmd)
}

func updateCLI() error {
	// Get current CLI version
	currentVersion := litetable.CLIVersion // Assuming Version is defined in your root.go or similar

	fmt.Printf("Current CLI version: %s\n", currentVersion)
	fmt.Println("Checking for updates...")

	// Get latest CLI version from GitHub API
	latestVersion, err := getLatestCLIVersion()
	if err != nil {
		return fmt.Errorf("failed to check for latest version: %w", err)
	}

	// Compare versions
	if !isUpdateAvailable(currentVersion, latestVersion) {
		fmt.Println("✅ Your CLI is already running the latest version!")
		return nil
	}

	fmt.Printf("New CLI version available: %s → %s\n", currentVersion, latestVersion)

	// Add confirmation prompt
	fmt.Print("Do you want to update? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	response = strings.ToLower(strings.TrimSpace(response))
	if response != "y" && response != "yes" {
		fmt.Println("Update canceled.")
		return nil
	}

	// Run the install script
	fmt.Println("Starting update process...")

	// Create install command
	var cmd *exec.Cmd

	// Execute the installation script with version parameter
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-c",
			fmt.Sprintf("iwr -useb https://raw.githubusercontent.com/litetable/litetable-cli/main/install.sh | bash -s %s", latestVersion))
	} else {
		cmd = exec.Command("bash", "-c",
			fmt.Sprintf("curl -fsSL https://raw.githubusercontent.com/litetable/litetable-cli/main/install.sh | bash -s %s", latestVersion))
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	fmt.Printf("✅  Successfully updated CLI to version %s!\n", latestVersion)
	return nil
}

func getLatestCLIVersion() (string, error) {
	// Use the GitHub API to get the latest release version
	cmd := exec.Command("curl", "-s", "https://api.github.com/repos/litetable/litetable-cli/releases/latest")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to fetch latest release info: %w", err)
	}

	// Simple regex to extract the tag_name from JSON response
	// For a more robust solution, consider using a JSON parser
	re := regexp.MustCompile(`"tag_name":\s*"(v\d+\.\d+\.\d+)"`)
	matches := re.FindSubmatch(output)
	if len(matches) < 2 {
		return "", fmt.Errorf("couldn't parse version from GitHub response")
	}

	return string(matches[1]), nil
}

// Reuse the version comparison logic from service/update.go
func isUpdateAvailable(currentVersion, latestVersion string) bool {
	return strings.Compare(latestVersion, currentVersion) > 0
}
