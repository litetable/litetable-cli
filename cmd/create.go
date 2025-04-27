package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"net/url"
	"strings"
	"time"
)

var (
	// Create command options
	families string

	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create configuration in the Litetable server",
		Long:  "Create configuration elements such as column families in the Litetable server",
		Run: func(cmd *cobra.Command, args []string) {
			if err := createConfig(); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		},
	}
)

func init() {
	// Add create command to root command
	rootCmd.AddCommand(createCmd)

	// Add flags for create operation
	createCmd.Flags().StringVarP(&families, "family", "f", "", "Column families to create (comma-separated)")
	_ = createCmd.MarkFlagRequired("family")
}

func createConfig() error {
	conn, err := dial()
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}
	defer conn.Close()

	now := time.Now()

	// Parse comma-separated families
	familyList := strings.Split(families, ",")
	for i, family := range familyList {
		familyList[i] = strings.TrimSpace(family)
	}

	// Encode family names to handle special characters
	encodedFamilies := make([]string, len(familyList))
	for i, family := range familyList {
		encodedFamilies[i] = url.QueryEscape(family)
	}

	// Create the CREATE command with all families
	cmd := fmt.Sprintf("CREATE family=%s", strings.Join(encodedFamilies, ","))

	// Send the command
	if _, err = conn.Write([]byte(cmd)); err != nil {
		return fmt.Errorf("failed to send create command: %w", err)
	}

	// Read response
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Try to parse as JSON, but handle non-JSON responses as well
	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(buffer[:n], &jsonResponse); err != nil {
		// Not JSON, just print the raw response
		fmt.Println(string(buffer[:n]))
	} else {
		// JSON response
		prettyJSON, _ := json.MarshalIndent(jsonResponse, "", "  ")
		fmt.Println(string(prettyJSON))
	}

	elapsed := time.Since(now)
	elapsedMs := float64(elapsed.Nanoseconds()) / 1_000_000.0
	fmt.Printf("Roundtrip in %.2fms\n", elapsedMs)

	return nil
}
