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
	"regexp"
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
	currentVersion, err := getInstalledVersion()
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
	if !isUpdateAvailable(currentVersion, latestVersion) {
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

	// Get LiteTable directory
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

	// Update config file with new version
	configPath := filepath.Join(liteTableDir, "litetable.conf")
	if err := updateConfigVersion(configPath, latestVersion); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	fmt.Printf("✅ Successfully updated to LiteTable version %s!\n", latestVersion)
	return nil
}

func getInstalledVersion() (string, error) {
	// Get LiteTable directory
	liteTableDir, err := dir.GetLitetableDir()
	if err != nil {
		return "", fmt.Errorf("failed to get LiteTable directory: %w", err)
	}

	// Read the config file
	configPath := filepath.Join(liteTableDir, "litetable.conf")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return "", fmt.Errorf("LiteTable is not installed or configuration file not found")
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse the config file to find version
	configLines := strings.Split(string(configBytes), "\n")
	for _, line := range configLines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, VersionKey) {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				version := strings.TrimSpace(parts[1])
				if version != "" {
					return version, nil
				}
			}
		}
	}

	return "", fmt.Errorf("version information not found in configuration")
}

func getLatestVersion() (string, error) {
	// Use git to list remote tags and get the latest version
	cmd := exec.Command("git", "ls-remote", "--tags", serverRepo)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to fetch remote tags: %w", err)
	}

	// Parse output to find the latest version tag
	re := regexp.MustCompile(`refs/tags/(v\d+\.\d+\.\d+)$`)
	var versions []string

	for _, line := range strings.Split(string(output), "\n") {
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			versions = append(versions, matches[1])
		}
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("no version tags found")
	}

	// Sort versions and return the latest
	// For simplicity, we'll rely on string comparison which works for vX.Y.Z format
	latestVersion := versions[0]
	for _, v := range versions[1:] {
		if strings.Compare(v, latestVersion) > 0 {
			latestVersion = v
		}
	}

	return latestVersion, nil
}

func isUpdateAvailable(currentVersion, latestVersion string) bool {
	return strings.Compare(latestVersion, currentVersion) > 0
}

func isServerRunning(liteTableDir string) bool {
	pidFile := filepath.Join(liteTableDir, "litetable.pid")
	_, err := os.Stat(pidFile)
	return err == nil
}

func updateConfigVersion(configPath, newVersion string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	configLines := strings.Split(string(data), "\n")
	updated := false

	for i, line := range configLines {
		if strings.HasPrefix(strings.TrimSpace(line), VersionKey) {
			configLines[i] = fmt.Sprintf("%s = %s", VersionKey, newVersion)
			updated = true
			break
		}
	}

	if !updated {
		configLines = append(configLines, fmt.Sprintf("%s = %s", VersionKey, newVersion))
	}

	return os.WriteFile(configPath, []byte(strings.Join(configLines, "\n")), 0644)
}
