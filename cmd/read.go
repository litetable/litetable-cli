package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/litetable/litetable-cli/internal/dir"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
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
	certDir, err := dir.GetLitetableDir()
	if err != nil {
		return fmt.Errorf("failed to get Litetable directory: %w", err)
	}

	// Path to the certificate file
	certFile := filepath.Join(certDir, serverCertName)

	// Load the server's certificate to trust it
	certData, err := os.ReadFile(certFile)
	if err != nil {
		return fmt.Errorf("failed to read certificate: %w", err)
	}

	// Create a certificate pool and add the server certificate
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(certData); !ok {
		return fmt.Errorf("failed to append certificate to pool")
	}

	// Create a TLS configuration that trusts the server certificate
	tlsConfig := &tls.Config{
		RootCAs:    certPool,
		ServerName: "localhost",
	}

	// Connect to the server using TLS
	conn, err := tls.Dial("tcp", ":9443", tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer conn.Close()

	fmt.Println("Connected to Litetable server")

	// Send some data
	message := []byte("Hello, Litetable Server!")
	_, err = conn.Write(message)
	if err != nil {
		return fmt.Errorf("failed to send data: %w", err)
	}

	// Read response
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	fmt.Printf("Received response: %s\n", buffer[:n])
	return nil
}
