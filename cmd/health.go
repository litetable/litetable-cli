package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net"
)

// healthCmd represents the health command
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check if the LiteTable server is running",
	Long:  `Sends a PING request to the LiteTable server to verify it's operational.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return checkServerHealth()
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)
}

func checkServerHealth() error {
	fmt.Println("🔍 Checking LiteTable server health...")

	conn, err := dial()
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
	cmd := fmt.Sprintf("PING ")

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

	fmt.Println(string(buffer[:n]))

	return nil
}
