package operations

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/litetable/litetable-cli/cmd/service"
	"github.com/litetable/litetable-cli/internal/litetable"
	"github.com/spf13/cobra"
	"io"
	"time"
)

var (
	// Read command options
	readKey       string
	readKeyPrefix string
	readRegex     string
	readFamily    string
	readQualifier []string
	readLatest    int

	ReadCmd = &cobra.Command{
		Use:   "read",
		Short: "Read data from the Litetable server",
		Long:  "Read allows you to retrieve data from the Litetable server",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Count how many key selectors are provided
			selectors := 0
			if readKey != "" {
				selectors++
			}
			if readKeyPrefix != "" {
				selectors++
			}
			if readRegex != "" {
				selectors++
			}

			if selectors != 1 {
				return fmt.Errorf("exactly one of --key (-k), --keyPrefix (-p), or --regex (-r) must be provided")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := readData(); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		},
	}
)

func init() {
	// Add flags for read operation
	ReadCmd.Flags().StringVarP(&readKey, "key", "k", "", "Row key to read")
	ReadCmd.Flags().StringVarP(&readKeyPrefix, "keyPrefix", "p", "",
		"Read all row-keys with this prefix")
	ReadCmd.Flags().StringVarP(&readRegex, "regex", "r", "",
		"Read all row-keys matching this regex pattern")
	ReadCmd.Flags().StringVarP(&readFamily, "family", "f", "", "Column family to read")
	ReadCmd.Flags().StringArrayVarP(&readQualifier, "qualifier", "q", []string{}, "Qualifiers to read (can be specified multiple times)")
	ReadCmd.Flags().IntVarP(&readLatest, "latest", "l", 0, "Number of latest versions to return")

	// Mark required flags - removing the required mark for key
	_ = ReadCmd.MarkFlagRequired("family")
}

func readData() error {
	conn, err := service.Dial()
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}
	defer conn.Close()

	now := time.Now()

	// Build the READ command based on which selector is provided
	var command string
	if readKey != "" {
		command = fmt.Sprintf("READ key=%s", readKey)
	} else if readKeyPrefix != "" {
		command = fmt.Sprintf("READ prefix=%s", readKeyPrefix)
	} else if readRegex != "" {
		// Create a properly formatted regex pattern that escapes special characters
		// and wraps the user input with ".*" for substring matching
		formattedRegex := fmt.Sprintf(".*%s.*", readRegex)
		command = fmt.Sprintf("READ regex=%s", formattedRegex)
	}

	if readFamily != "" {
		command += fmt.Sprintf(" family=%s", readFamily)
	}

	for _, q := range readQualifier {
		command += fmt.Sprintf(" qualifier=%s", q)
	}

	if readLatest > 0 {
		command += fmt.Sprintf(" latest=%d", readLatest)
	}

	// Send the command
	if _, err = conn.Write([]byte(command)); err != nil {
		return fmt.Errorf("failed to send read command: %w", err)
	}

	// Read response using the robust approach for large responses
	var fullResponse []byte
	buffer := make([]byte, 4096)

	// Use a reasonable timeout for the entire read operation
	timeout := time.Now().Add(10 * time.Second)
	if err := conn.SetReadDeadline(timeout); err != nil {
		return fmt.Errorf("failed to set read deadline: %w", err)
	}

	for {
		n, err := conn.Read(buffer)
		if n > 0 {
			fullResponse = append(fullResponse, buffer[:n]...)

			// Check if we have a complete JSON object
			if len(fullResponse) > 0 && IsValidJSON(fullResponse) {
				break
			}
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("error reading response: %w", err)
		}

		// Extend deadline for each successful read
		if err := conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
			return fmt.Errorf("failed to extend read deadline: %w", err)
		}
	}

	elapsed := time.Since(now)
	elapsedMs := float64(elapsed.Nanoseconds()) / 1_000_000.0

	// Parse the response based on which query mode we're using
	// Parse as an array of rows
	var rows map[string]litetable.Row
	if err := json.Unmarshal(fullResponse, &rows); err != nil {
		return fmt.Errorf("%s", string(fullResponse))
	}

	// Print each row with a separator
	first := true
	for key, row := range rows {
		if !first {
			fmt.Println("--------------------")
		}
		first = false
		fmt.Printf("Key: %s\n%s\n", key, row.PrettyPrint())
	}
	fmt.Printf("Row results: %d\n", len(rows))
	fmt.Printf("Response size: %d bytes\n", len(fullResponse))
	fmt.Printf("Query duration: %.2fms\n", elapsedMs)

	return nil
}

// IsValidJSON checks if the buffer contains a complete, valid JSON object
func IsValidJSON(data []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(data, &js) == nil
}
