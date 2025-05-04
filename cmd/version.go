package cmd

import (
	"fmt"
	"github.com/litetable/litetable-cli/cmd/service"
	"github.com/litetable/litetable-cli/internal/dir"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Display LiteTable version information",
	Long:  "Show the installed version of LiteTable server from configuration",
	Run: func(cmd *cobra.Command, args []string) {
		version, err := getInstalledVersion()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("LiteTable Server Version: %s\n", version)
	},
}

func init() {
	// Add the version command to the root command
	rootCmd.AddCommand(versionCommand)
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
		if strings.HasPrefix(line, service.ServiceVersionKey) {
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
