package cmd

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/litetable/litetable-cli/internal/dir"
	"net"
	"os"
	"path/filepath"
)

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
