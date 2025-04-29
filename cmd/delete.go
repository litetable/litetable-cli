package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"time"
)

var (
	// Delete command options
	deleteKey       string
	deleteFamily    string
	deleteQualifier []string
	deleteTTL       int64
	deleteFrom      string

	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete data from the Litetable server",
		Long:  "Delete allows you to remove data from the Litetable server",
		Run: func(cmd *cobra.Command, args []string) {
			// Execute the delete operation
			if err := deleteData(); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		},
	}
)

func init() {
	// Register delete command with root command
	rootCmd.AddCommand(deleteCmd)

	// Add flags for delete operation
	deleteCmd.Flags().StringVarP(&deleteKey, "key", "k", "", "Row key to delete (required)")
	deleteCmd.Flags().StringVarP(&deleteFamily, "family", "f", "", "Column family to delete")
	deleteCmd.Flags().StringArrayVarP(&deleteQualifier, "qualifier", "q", []string{}, "Qualifiers to delete (can be specified multiple times)")
	deleteCmd.Flags().Int64Var(&deleteTTL, "ttl", 0, "Time-to-live in seconds for tombstone entries")
	deleteCmd.Flags().StringVar(&deleteFrom, "from", "", "Starting position for deletion in the map")

	// Mark required flags
	_ = deleteCmd.MarkFlagRequired("key")
}

func deleteData() error {
	conn, err := dial()
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}
	defer conn.Close()

	now := time.Now()

	// Build the DELETE command
	command := fmt.Sprintf("DELETE key=%s", deleteKey)

	if deleteFamily != "" {
		command += fmt.Sprintf(" family=%s", deleteFamily)
	}

	for _, q := range deleteQualifier {
		command += fmt.Sprintf(" qualifier=%s", q)
	}

	// Add the optional from parameter if provided
	if deleteFrom != "" {
		command += fmt.Sprintf(" timestamp=%s", deleteFrom)
	}

	// Add the optional TTL parameter if provided
	if deleteTTL > 0 {
		command += fmt.Sprintf(" ttl=%d", deleteTTL)
	}

	// Send the command
	if _, err = conn.Write([]byte(command)); err != nil {
		return fmt.Errorf("failed to send delete command: %w", err)
	}

	// Read response using a more robust approach
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
			if len(fullResponse) > 0 && isValidJSON(fullResponse) {
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

	// Print raw response size for debugging
	fmt.Printf("Received %d bytes\n", len(fullResponse))

	// Parse the response
	elapsed := time.Since(now)
	elapsedMs := float64(elapsed.Nanoseconds()) / 1_000_000.0
	fmt.Printf("Deletion completed in %.2fms\n", elapsedMs)

	// Pretty print the response
	fmt.Printf("Result: %v\n", string(fullResponse))
	return nil
}
