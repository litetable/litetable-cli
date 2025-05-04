package litetable

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// GetLatestVersion fetches the latest version of a provided git repository URL
func GetLatestVersion(url string) (string, error) {
	// Use git to list remote tags and get the latest version
	cmd := exec.Command("git", "ls-remote", "--tags", url)
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
