package cmd

import (
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/litetable/litetable-cli/internal/dir"
	"github.com/litetable/litetable-cli/internal/litetable"
	"github.com/spf13/cobra"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"
)

var (
	// Read command options
	readKey       string
	readFamily    string
	readQualifier []string
	readLatest    int

	readCmd = &cobra.Command{
		Use:   "read",
		Short: "Read data from the Litetable server",
		Long:  "Read allows you to retrieve data from the Litetable server",
		Run: func(cmd *cobra.Command, args []string) {
			if err := readData(); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		},
	}
)

func init() {
	// Add read command to root command
	rootCmd.AddCommand(readCmd)

	// Add flags for read operation
	readCmd.Flags().StringVarP(&readKey, "key", "k", "", "Row key to read (required)")
	readCmd.Flags().StringVarP(&readFamily, "family", "f", "", "Column family to read")
	readCmd.Flags().StringArrayVarP(&readQualifier, "qualifier", "q", []string{}, "Qualifiers to read (can be specified multiple times)")
	readCmd.Flags().IntVarP(&readLatest, "latest", "l", 0, "Number of latest versions to return")

	// Mark required flags
	_ = readCmd.MarkFlagRequired("key")
	_ = readCmd.MarkFlagRequired("family")
}

func readData() error {
	conn, err := dial()
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}
	defer conn.Close()

	now := time.Now()

	// Build the READ command
	command := fmt.Sprintf("READ key=%s", readKey)

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

	var payload litetable.Row
	if err := json.Unmarshal(fullResponse, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w\nRaw: %s",
			err, string(fullResponse))
	}

	elapsed := time.Since(now)
	elapsedMs := float64(elapsed.Nanoseconds()) / 1_000_000.0
	fmt.Printf("Roundtrip in %.2fms\n\n%s\n", elapsedMs, payload.PrettyPrint())
	return nil
}

func dial() (net.Conn, error) {
	certDir, err := dir.GetLitetableDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get Litetable directory: %w", err)
	}

	// Path to the certificate file
	certFile := filepath.Join(certDir, serverCertName)

	// Load the server's certificate to trust it
	certData, err := os.ReadFile(certFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate: %w", err)
	}

	// Create a certificate pool and add the server certificate
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(certData); !ok {
		return nil, fmt.Errorf("failed to append certificate to pool")
	}

	// Create a TLS configuration that trusts the server certificate
	// tlsConfig := &tls.Config{
	// 	RootCAs:            certPool,
	// 	ServerName:         "localhost",
	// }

	// Connect to the server using TLS
	conn, err := net.Dial("tcp", ":9443")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	return conn, nil
}

// isValidJSON checks if the buffer contains a complete, valid JSON object
func isValidJSON(data []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(data, &js) == nil
}
