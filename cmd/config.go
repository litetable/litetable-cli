package cmd

import (
	"fmt"
	"github.com/litetable/litetable-cli/internal/dir"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	cfgKey                string
	cfgVal                string
	allowedConfigurations = []string{
		"server_port",
		"server_address",
		"garbage_collection_timer",
		"backup_timer",
		"snapshot_timer",
		"max_snapshot_limit",
		"debug",
		"cloud_environment",
	}
)

// Config represents the structure of the config file
type Config map[string]string

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(setCmd)
	configCmd.AddCommand(updateCmd)

	// Add common flags for set and update commands
	setCmd.Flags().StringVarP(&cfgKey, "key", "k", "", "Config key to set (required)")
	setCmd.Flags().StringVarP(&cfgVal, "value", "v", "", "Config value to set (required)")
	setCmd.MarkFlagRequired("key")
	setCmd.MarkFlagRequired("value")

	updateCmd.Flags().StringVarP(&cfgKey, "key", "k", "", "Config key to update (required)")
	updateCmd.Flags().StringVarP(&cfgVal, "value", "v", "", "Config value to update (required)")
	updateCmd.MarkFlagRequired("key")
	updateCmd.MarkFlagRequired("value")

	configCmd.AddCommand(viewCmd)

}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage litetable configuration",
	Long:  `Configure litetable settings through set and update operations.`,
}

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View the current configuration",
	Long:  `Display the contents of the litetable.conf file.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, err := getConfigFilePath()
		if err != nil {
			fmt.Printf("Error finding config file: %s\n", err)
			os.Exit(1)
		}

		// Check if the file exists
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			fmt.Println("Configuration file does not exist.")
			return
		}

		// Read the file contents
		data, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Printf("Error reading config file: %s\n", err)
			os.Exit(1)
		}

		if len(data) == 0 {
			fmt.Println("Configuration file is empty.")
			return
		}

		// Print the raw content
		fmt.Println(string(data))
	},
}

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a configuration value",
	Long:  `Set a new configuration value. Will error if the key already exists.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cfgKey == "" {
			fmt.Println("Error: key is required")
			return
		}
		if cfgVal == "" {
			fmt.Println("Error: value is required")
			return
		}

		// Validate if key is allowed
		if !isAllowedConfigKey(cfgKey) {
			fmt.Printf("Error: '%s' is not an allowed configuration key.\n", cfgKey)
			fmt.Printf("Allowed keys: %s\n", strings.Join(allowedConfigurations, ", "))
			return
		}

		// Load config
		config, err := loadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %s\n", err)
			return
		}

		// Check if key already exists
		if _, exists := config[cfgKey]; exists {
			fmt.Printf("Error: key '%s' already exists. Use 'update' to modify existing values.\n", cfgKey)
			return
		}

		// Set the value
		config[cfgKey] = cfgVal
		if err := saveConfig(config); err != nil {
			fmt.Printf("Error saving config: %s\n", err)
			return
		}
		fmt.Printf("Successfully set '%s' to '%s'\n", cfgKey, cfgVal)
	},
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing configuration value",
	Long:  `Update an existing configuration value. Will error if the key does not exist.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cfgKey == "" {
			fmt.Println("Error: key is required")
			return
		}
		if cfgVal == "" {
			fmt.Println("Error: value is required")
			return
		}

		// Validate if key is allowed
		if !isAllowedConfigKey(cfgKey) {
			fmt.Printf("Error: '%s' is not an allowed configuration key.\n", cfgKey)
			fmt.Printf("Allowed keys: %s\n", strings.Join(allowedConfigurations, ", "))
			return
		}

		// Load config
		config, err := loadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %s\n", err)
			return
		}

		// Check if key exists
		if _, exists := config[cfgKey]; !exists {
			fmt.Printf("Error: key '%s' does not exist. Use 'set' to create a new value.\n", cfgKey)
			return
		}

		// Update the value
		config[cfgKey] = cfgVal
		if err := saveConfig(config); err != nil {
			fmt.Printf("Error updating config: %s\n", err)
			return
		}
		fmt.Printf("Successfully updated '%s' to '%s'\n", cfgKey, cfgVal)
	},
}

// getConfigFilePath returns the path to the config file
func getConfigFilePath() (string, error) {
	home, err := dir.GetLitetableDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(home, "litetable.conf"), nil
}

// loadConfig reads the configuration from the config file
func loadConfig() (Config, error) {
	config := make(Config)

	configPath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	// If file doesn't exist, return empty config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Empty file case
	if len(data) == 0 {
		return config, nil
	}

	// Parse the config line by line
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		// Skip comments and empty lines
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split by = sign
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		config[key] = value
	}

	return config, nil
}

// saveConfig writes the configuration to the config file
func saveConfig(config Config) error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// Read existing file to preserve comments
	var lines []string
	existingConfig := make(map[string]bool)
	fileExists := false

	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		fileExists = true
		fileData, err := os.ReadFile(configPath)
		if err == nil {
			fileLines := strings.Split(string(fileData), "\n")
			for _, line := range fileLines {
				// Skip empty lines at the end of file
				if len(lines) > 0 && line == "" && len(fileLines) > 1 {
					continue
				}

				trimmedLine := strings.TrimSpace(line)
				if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
					// Keep comments and empty lines
					lines = append(lines, line)
				} else {
					// Extract key from existing config lines
					parts := strings.SplitN(trimmedLine, "=", 2)
					if len(parts) == 2 {
						key := strings.TrimSpace(parts[0])
						existingConfig[key] = true

						// If key exists in the new config, update its value
						if val, ok := config[key]; ok {
							lines = append(lines, fmt.Sprintf("%s = %s", key, val))
							delete(config, key) // Remove from config to avoid duplication
						} else {
							lines = append(lines, line) // Keep line as is
						}
					} else {
						lines = append(lines, line) // Keep line as is
					}
				}
			}
		}
	}

	// Append any new keys with proper formatting
	if len(config) > 0 {
		// Add a blank line between existing content and new entries if file exists and doesn't end with blank line
		if fileExists && len(lines) > 0 && lines[len(lines)-1] != "" {
			lines = append(lines, "")
		}

		// Add new configurations
		for key, value := range config {
			if !existingConfig[key] {
				lines = append(lines, fmt.Sprintf("%s = %s", key, value))
			}
		}
	}

	// Ensure we end with a newline
	output := strings.Join(lines, "\n")
	if !strings.HasSuffix(output, "\n") && len(output) > 0 {
		output += "\n"
	}

	return os.WriteFile(configPath, []byte(output), 0644)
}

// isAllowedConfigKey checks if the provided key is in the allowed configurations list
func isAllowedConfigKey(key string) bool {
	for _, allowedKey := range allowedConfigurations {
		if key == allowedKey {
			return true
		}
	}
	return false
}
