package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/litetable/litetable-cli/internal/dir"
	"github.com/litetable/litetable-cli/internal/litetable"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
	"time"
)

var (
	readCmd = &cobra.Command{
		Use:   "read",
		Short: "Read data from the Litetable server",
		Long:  "Read allows you to retrieve data from the Litetable server",
		Run: func(cmd *cobra.Command, args []string) {
			// read data from a read query
			if err := readData(); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		},
	}
)

func init() {
	// add read command to root command
	rootCmd.AddCommand(readCmd)
}

// readData requires us to create a connection to the server before sending a protocol
// message over TLS. Connections require a TLS certificate.
func readData() error {
	conn, err := dial()
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}
	defer conn.Close()

	now := time.Now()

	// Send the command
	message := []byte("READ key=testKey:12345 family=main qualifier=status qualifier=time latest" +
		"=10")
	if _, err = conn.Write(message); err != nil {
		return fmt.Errorf("failed to send data: %w", err)
	}

	// Read response using a more robust approach for large responses
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

func dial() (*tls.Conn, error) {
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
	tlsConfig := &tls.Config{
		RootCAs:    certPool,
		ServerName: "localhost",
	}

	// Connect to the server using TLS
	conn, err := tls.Dial("tcp", ":9443", tlsConfig)
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
