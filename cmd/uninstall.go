package cmd

import (
	"bufio"
	"fmt"
	"github.com/litetable/litetable-cli/internal/dir"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var uninstallCommand = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall LiteTable CLI",
	Long:  "Removes the LiteTable CLI binary directory while preserving your data and configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if err := uninstallLiteTable(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func uninstallLiteTable() error {
	// Get LiteTable directory
	liteTableDir, err := dir.GetLitetableDir()
	if err != nil {
		return fmt.Errorf("failed to get LiteTable directory: %w", err)
	}

	// Path to bin directory
	binDir := filepath.Join(liteTableDir, "bin")

	// Check if bin directory exists
	if _, err := os.Stat(binDir); os.IsNotExist(err) {
		fmt.Println("LiteTable bin directory not found. Nothing to uninstall.")
		return nil
	}

	// Add confirmation prompt
	fmt.Print("Do you want to uninstall LiteTable? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	response = strings.ToLower(strings.TrimSpace(response))
	if response != "y" && response != "yes" {
		fmt.Println("Uninstall canceled.")
		return nil
	}

	// Start the update
	fmt.Println("Starting update process...")

	fmt.Println("Uninstalling LiteTable CLI...")
	fmt.Printf("Removing binary directory: %s\n", binDir)

	// Remove bin directory
	if err := os.RemoveAll(binDir); err != nil {
		return fmt.Errorf("failed to remove bin directory: %w", err)
	}

	// Remove LiteTable from PATH in shell config files
	if err := removeFromPath(liteTableDir); err != nil {
		fmt.Printf("Warning: Failed to remove LiteTable from PATH: %v\n", err)
		fmt.Println("You may need to manually remove the LiteTable PATH entry from your shell configuration file.")
	}

	fmt.Println("✅ LiteTable CLI binaries have been successfully removed.")
	fmt.Println("Your configuration and data are preserved.")
	fmt.Println("To completely remove LiteTable, delete the directory:", liteTableDir)

	return nil
}

// removeFromPath attempts to remove the LiteTable bin directory from the PATH in various shell config files
func removeFromPath(liteTableDir string) error {
	if runtime.GOOS == "windows" {
		fmt.Println("Note: On Windows, you may need to manually update your PATH environment variable.")
		return nil
	}

	binDir := filepath.Join(liteTableDir, "bin")
	pathPattern := fmt.Sprintf("export PATH=.*%s.*", regexp.QuoteMeta(binDir))

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Common shell config files
	configFiles := []string{
		filepath.Join(homeDir, ".bashrc"),
		filepath.Join(homeDir, ".bash_profile"),
		filepath.Join(homeDir, ".zshrc"),
		filepath.Join(homeDir, ".profile"),
	}

	for _, configFile := range configFiles {
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			continue
		}

		data, err := os.ReadFile(configFile)
		if err != nil {
			continue
		}

		content := string(data)
		re := regexp.MustCompile(pathPattern)
		if !re.MatchString(content) {
			continue
		}

		// Remove the PATH entry
		newContent := re.ReplaceAllString(content, "")

		// Remove any empty lines that might have been created
		emptyLineRe := regexp.MustCompile(`\n\n+`)
		newContent = emptyLineRe.ReplaceAllString(newContent, "\n\n")

		if err := os.WriteFile(configFile, []byte(newContent), 0644); err != nil {
			return fmt.Errorf("failed to update %s: %w", configFile, err)
		}

		fmt.Printf("✓ Removed LiteTable from PATH in %s\n", configFile)
	}

	return nil
}
