package cmd

import (
	"bufio"
	"fmt"
	"github.com/litetable/litetable-cli/internal/litetable"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	installUrl = litetable.CLIInstallURL
	UpdateCmd  = &cobra.Command{
		Use:   "update",
		Short: "Update LiteTable CLI",
		Long:  "Check for and install the latest version of LiteTable CLI",
		Run: func(cmd *cobra.Command, args []string) {
			if err := updateCLI(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(UpdateCmd)
}

func updateCLI() error {
	// Get current CLI version
	currentVersion, err := litetable.GetFromConfig(litetable.CLIVersionKey)
	if err != nil {
		return fmt.Errorf("failed to get current CLI version: %w", err)
	}

	fmt.Printf("Current CLI version: %s\nChecking for updates...\n", currentVersion)

	// Get the latest CLI version from GitHub API
	latestVersion, err := litetable.GetLatestVersion(litetable.CLIURL)
	if err != nil {
		return fmt.Errorf("failed to check for latest version: %w", err)
	}

	// Compare versions
	if !litetable.IsUpdateAvailable(currentVersion, latestVersion) {
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

	// Start the update
	fmt.Println("Starting update process...")

	// Create install command
	var cmd *exec.Cmd

	// Execute the installation script with version parameter
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-c",
			fmt.Sprintf("iwr -useb %s | bash -s %s", installUrl, latestVersion))
	} else {
		cmd = exec.Command("bash", "-c",
			fmt.Sprintf("curl -fsSL %s | bash -s %s", installUrl, latestVersion))
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	// update the version to latest
	if err = litetable.UpdateConfigValue(&litetable.UpdateConfig{
		Key:   litetable.CLIVersionKey,
		Value: latestVersion,
	}); err != nil {
		return fmt.Errorf("failed to update version in config: %w", err)
	}
	fmt.Printf("✅  Successfully updated CLI to version %s!\n", latestVersion)
	return nil
}
