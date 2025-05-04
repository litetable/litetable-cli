package litetable

import (
	"fmt"
	"github.com/litetable/litetable-cli/internal/dir"
	"os"
	"path/filepath"
	"strings"
)

const (
	ServerVersionKey = "server_version"
)

func GetFromConfig(value string) (string, error) {
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
		if strings.HasPrefix(line, value) {
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

type UpdateConfig struct {
	Key   string
	Value string
}

// UpdateConfigValue updates the LiteTable configuration file with the provided key and value
func UpdateConfigValue(cfg *UpdateConfig) error {
	// Get LiteTable directory
	liteTableDir, err := dir.GetLitetableDir()
	if err != nil {
		return fmt.Errorf("failed to get LiteTable directory: %w", err)
	}

	// Update the config file with provided values
	configPath := filepath.Join(liteTableDir, "litetable.conf")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	configLines := strings.Split(string(data), "\n")
	updated := false

	for i, line := range configLines {
		if strings.HasPrefix(strings.TrimSpace(line), cfg.Key) {
			configLines[i] = fmt.Sprintf("%s = %s", cfg.Key, cfg.Value)
			updated = true
			break
		}
	}

	if !updated {
		configLines = append(configLines, fmt.Sprintf("%s = %s", cfg.Key, cfg.Value))
	}

	return os.WriteFile(configPath, []byte(strings.Join(configLines, "\n")), 0644)
}
