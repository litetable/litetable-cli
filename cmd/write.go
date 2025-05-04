package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/litetable/litetable-cli/cmd/service"
	"github.com/litetable/litetable-cli/internal/litetable"
	"github.com/spf13/cobra"
	"net"
	"net/url"
	"strconv"
	"time"
)

var (
	writeKey    string
	writeFamily string
	writeQuals  []string
	writeValues []string
	writeTTL    int64

	writeCmd = &cobra.Command{
		Use:   "write",
		Short: "Write data to the Litetable server",
		Long:  "Write allows you to send data to the Litetable server",
		Example: "litetable write --key=rowKey --family=familyName --qualifier=qual1 --value=val1" +
			" --qualifier=qual2 --value=val2 --ttl=60",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Validate inputs
			if writeKey == "" {
				return fmt.Errorf("key is required")
			}
			if writeFamily == "" {
				return fmt.Errorf("family is required")
			}
			if len(writeQuals) != len(writeValues) {
				return fmt.Errorf("number of qualifiers (%d) must match number of values (%d)",
					len(writeQuals), len(writeValues))
			}
			if len(writeQuals) == 0 {
				return fmt.Errorf("at least one qualifier/value pair is required")
			}
			if writeTTL < 0 {
				return fmt.Errorf("TTL must be a non-negative value")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := writeData(); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		},
	}
)

func init() {
	writeCmd.Flags().StringVarP(&writeKey, "key", "k", "", "Row key for the write operation")
	writeCmd.Flags().StringVarP(&writeFamily, "family", "f", "",
		"Column family for the write operation")
	writeCmd.Flags().StringArrayVarP(&writeQuals, "qualifier", "q", []string{},
		"Qualifiers to read (can be specified multiple times)")
	writeCmd.Flags().StringArrayVarP(&writeValues, "value", "v", []string{},
		"Values to write (can be specified multiple times, use quotes for values with spaces)")
	writeCmd.Flags().Int64VarP(&writeTTL, "ttl", "t", 0,
		"Time to live in seconds (0 means no expiration)")

	rootCmd.AddCommand(writeCmd)
}

func writeData() error {
	conn, err := service.Dial()
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}

	defer func(conn net.Conn) {
		closeErr := conn.Close()
		if closeErr != nil {
			fmt.Println("failed to close connection: %w", closeErr.Error())
		}
	}(conn)

	// Create the WRITE command with all the qualifier/value pairs
	cmd := fmt.Sprintf("WRITE key=%s family=%s", writeKey, writeFamily)
	for i := 0; i < len(writeQuals); i++ {
		// URL encode the value to properly handle spaces and special characters
		encodedValue := url.QueryEscape(writeValues[i])
		cmd += fmt.Sprintf(" qualifier=%s value=%s", writeQuals[i], encodedValue)
	}

	if writeTTL > 0 {
		cmd += fmt.Sprintf(" ttl=%s", strconv.FormatInt(writeTTL, 10))
	}

	now := time.Now()
	// Send the write command
	_, err = conn.Write([]byte(cmd))
	if err != nil {
		return fmt.Errorf("failed to send data: %w", err)
	}

	// Read response
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var payload litetable.Row
	if err := json.Unmarshal(buffer[:n], &payload); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w\nRaw: %s",
			err, string(buffer[:n]))
	}
	elapsed := time.Since(now)
	elapsedMs := float64(elapsed.Nanoseconds()) / 1_000_000.0
	fmt.Printf("Roundtrip in %.2fms\n\n%s\n", elapsedMs, payload.PrettyPrint())
	return nil

}
