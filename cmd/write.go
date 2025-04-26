package cmd

import (
	"crypto/tls"
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

var (
	writeKey    string
	writeFamily string
	writeQuals  []string
	writeValues []string

	writeCmd = &cobra.Command{
		Use:     "write",
		Short:   "Write data to the Litetable server",
		Long:    "Write allows you to send data to the Litetable server",
		Example: "litetable write --key=rowKey --family=familyName --qualifier=qual1 --value=val1 --qualifier=qual2 --value=val2",
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
	writeCmd.Flags().StringVar(&writeKey, "key", "", "Row key for the write operation")
	writeCmd.Flags().StringVar(&writeFamily, "family", "", "Column family for the write operation")
	writeCmd.Flags().StringArrayVar(&writeQuals, "qualifier", []string{}, "Column qualifier (can be specified multiple times)")
	writeCmd.Flags().StringArrayVar(&writeValues, "value", []string{}, "Value to write (can be specified multiple times)")

	rootCmd.AddCommand(writeCmd)
}

func writeData() error {
	conn, err := dial()
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}

	defer func(conn *tls.Conn) {
		closeErr := conn.Close()
		if closeErr != nil {
			fmt.Println("failed to close connection: %w", closeErr.Error())
		}
	}(conn)

	// Create the WRITE command with all the qualifier/value pairs
	cmd := fmt.Sprintf("WRITE key=%s family=%s", writeKey, writeFamily)
	for i := 0; i < len(writeQuals); i++ {
		cmd += fmt.Sprintf(" qualifier=%s value=%s", writeQuals[i], writeValues[i])
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

	elapsed := time.Since(now)
	elapsedMs := float64(elapsed.Nanoseconds()) / 1_000_000.0
	fmt.Printf("Roundtrip in %.2fms\nResponse - %s\n", elapsedMs, string(buffer[:n]))
	return nil
}
