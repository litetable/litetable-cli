package dir

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	litetableDir = ".litetable"
	familiesFile = "families.config.json"
)

// GetLitetableDir returns the path to the LiteTable directory in the user's home directory.
func GetLitetableDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	dir := filepath.Join(homeDir, litetableDir)

	return dir, nil
}

// GetFamilies reads the families file and returns a list of families.
func GetFamilies() ([]string, error) {
	ltDir, err := GetLitetableDir()
	if err != nil {
		return nil, err
	}

	familiesFilePath := filepath.Join(ltDir, familiesFile)
	if _, err := os.Stat(familiesFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("families file does not exist: %w", err)
	}

	// Read the families file
	fileContent, err := os.ReadFile(familiesFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read families file: %w", err)
	}

	// Unmarshal the JSON array
	var families []string
	if err := json.Unmarshal(fileContent, &families); err != nil {
		return nil, fmt.Errorf("failed to unmarshal families file: %w", err)
	}

	return families, nil
}
