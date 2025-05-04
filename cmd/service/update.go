package service

import (
	"bufio"
	"fmt"
	"github.com/litetable/litetable-cli/internal/dir"
	"github.com/litetable/litetable-cli/internal/litetable"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var UpdateCommand = &cobra.Command{
	Use:   "update",
	Short: "Update LiteTable server",
	Long:  "Check for and install the latest version of LiteTable server",
	Run: func(cmd *cobra.Command, args []string) {
		if err := updateLiteTable(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	},
}

func updateLiteTable() error {
	// Get current version
	currentVersion, err := litetable.GetFromConfig(litetable.ServerVersionKey)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	fmt.Printf("Current version: %s\n", currentVersion)
	fmt.Println("Checking for updates...")

	// Get the latest version from git repository
	latestVersion, err := litetable.GetLatestVersion(serverRepo)
	if err != nil {
		return fmt.Errorf("failed to check for latest version: %w", err)
	}

	// Compare versions
	if !litetable.IsUpdateAvailable(currentVersion, latestVersion) {
		fmt.Println("✅  You are already running the latest version!")
		return nil
	}

	fmt.Printf("New version available: %s → %s\n", currentVersion, latestVersion)

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
	fmt.Println("Starting update process...")

	// Get LiteTable directory to verify running status
	liteTableDir, err := dir.GetLitetableDir()
	if err != nil {
		return fmt.Errorf("failed to get LiteTable directory: %w", err)
	}

	// Stop server if running
	if isServerRunning(liteTableDir) {
		fmt.Println("Stopping running LiteTable server...")
		if err := stopLiteTable(); err != nil {
			return fmt.Errorf("failed to stop server: %w", err)
		}
	}

	// Create temporary directory for cloning
	tempDir, err := ioutil.TempDir("", "litetable-update")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Clone the repository
	fmt.Println("📥 Downloading latest LiteTable server...")
	cloneCmd := exec.Command("git", "clone", "--depth", "1", "-b", latestVersion, serverRepo, tempDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	if err := cloneCmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// Build the server binary
	fmt.Println("🔨 Building server...")
	binDir := filepath.Join(liteTableDir, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	binPath := filepath.Join(binDir, serverBin)
	if runtime.GOOS == "windows" {
		binPath += ".exe"
	}

	buildCmd := exec.Command("go", "build", "-o", binPath)
	buildCmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%s", runtime.GOOS),
		fmt.Sprintf("GOARCH=%s", runtime.GOARCH))
	buildCmd.Dir = tempDir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build server: %w", err)
	}

	// Make executable (for Unix systems)
	if runtime.GOOS != "windows" {
		if err := os.Chmod(binPath, 0755); err != nil {
			return fmt.Errorf("failed to make server executable: %w", err)
		}
	}

	if err = litetable.UpdateConfigValue(&litetable.UpdateConfig{
		ConfigName: litetable.ServerVersionKey,
		NewVersion: latestVersion,
	}); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	fmt.Printf("✅ Successfully updated to LiteTable version %s!\n", latestVersion)
	return nil
}

func isServerRunning(liteTableDir string) bool {
	pidFile := filepath.Join(liteTableDir, "litetable.pid")
	_, err := os.Stat(pidFile)
	return err == nil
}
