package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var (
	forceWipe bool

	wipeCmd = &cobra.Command{
		Use:   "wipe",
		Short: "Wipe all LiteTable data",
		Long: "Wipe removes all LiteTable data files from the WAL, Garbage Collector, " +
			"and data backups. Does not remove server configuration.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := wipeData(); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		},
	}
)

func init() {
	wipeCmd.Flags().BoolVarP(&forceWipe, "force", "f", false, "Skip confirmation prompt")
}

func wipeData() error {
	// Define paths to remove
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	liteTableDir := filepath.Join(homeDir, ".litetable")

	familiesFile := filepath.Join(liteTableDir, "families.config.json")
	reaperFile := filepath.Join(liteTableDir, ".reaper.gc.log")

	tableDir := filepath.Join(liteTableDir, ".table")

	walDir := filepath.Join(liteTableDir, "wal")
	walLogFile := filepath.Join(walDir, "wal.log")

	// Check if directories/files exist
	var pathsToRemove []string

	if _, err := os.Stat(tableDir); err == nil {
		pathsToRemove = append(pathsToRemove, tableDir)
	}

	if _, err := os.Stat(walLogFile); err == nil {
		pathsToRemove = append(pathsToRemove, walLogFile)
	}

	if _, err := os.Stat(familiesFile); err == nil {
		pathsToRemove = append(pathsToRemove, familiesFile)
	}

	if _, err := os.Stat(reaperFile); err == nil {
		pathsToRemove = append(pathsToRemove, reaperFile)
	}

	if len(pathsToRemove) == 0 {
		fmt.Println("No LiteTable data found to wipe.")
		return nil
	}

	// Confirm deletion
	if !forceWipe {
		warningMsg := `
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃    WARNING: DESTRUCTIVE OPERATION    ┃
┃                                      ┃
┃  ALL LITETABLE DATA WILL BE DELETED  ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
`
		fmt.Println(warningMsg)
		fmt.Println("The following paths will be removed:")
		for _, path := range pathsToRemove {
			fmt.Printf("  - %s\n", path)
		}
		fmt.Print("\nTo confirm deletion, type 'DELETE' and press Enter: ")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read confirmation: %w", err)
		}

		response = strings.TrimSpace(response)
		if response != "DELETE" {
			fmt.Println("Operation canceled.")
			return nil
		}
	}

	// Delete the paths (handling both files and directories)
	for _, path := range pathsToRemove {
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("failed to access %s: %w", path, err)
		}

		if info.IsDir() {
			// Recursively remove directory and contents
			if err := os.RemoveAll(path); err != nil {
				return fmt.Errorf("failed to delete directory %s: %w", path, err)
			}
			fmt.Printf("Deleted directory: %s\n", path)
		} else {
			// Remove individual file
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to delete file %s: %w", path, err)
			}
			fmt.Printf("Deleted file: %s\n", path)
		}
	}

	// Try to remove empty directories
	if _, err := os.Stat(walDir); err == nil {
		if entries, err := os.ReadDir(walDir); err == nil && len(entries) == 0 {
			if err := os.Remove(walDir); err == nil {
				fmt.Printf("Removed empty directory: %s\n", walDir)
			}
		}
	}

	if _, err := os.Stat(liteTableDir); err == nil {
		if entries, err := os.ReadDir(liteTableDir); err == nil && len(entries) == 0 {
			if err := os.Remove(liteTableDir); err == nil {
				fmt.Printf("Removed empty directory: %s\n", liteTableDir)
			}
		}
	}

	fmt.Println("LiteTable data has been successfully wiped.")
	return nil
}
