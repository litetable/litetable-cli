package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/litetable/litetable-cli/internal/dir"
	"github.com/spf13/cobra"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

const (
	serverCertName = "server.crt"
	serverKeyName  = "server.key"
)

var (
	credentialsCmd = &cobra.Command{
		Use:   "credentials",
		Short: "Manage credentials",
		Long:  "Manage credentials for a Litetable Instance",
		Run: func(cmd *cobra.Command, args []string) {
			// default is run help
			_ = cmd.Help()
		},
	}

	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate credentials",
		Long:  "Generate credentials for a Litetable Instance",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Generating credentials...")
			// Your generation logic here
			if err := generateCredentials(); err != nil {
				fmt.Printf("Error generating credentials: %v\n", err)
				return
			}
		},
	}
)

func init() {
	// create a generate flag
	credentialsCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(credentialsCmd)
}

// generateCredentials generates an openSSL certificate that all instances
// of Litetable DB will use for TLS communication
func generateCredentials() error {
	// Create a certificate in the user's OS home directory in a private
	// directory called .litetable
	certDir, err := dir.GetLitetableDir()
	if err != nil {
		return fmt.Errorf("failed to get litetable directory: %w", err)
	}

	// create the directory if it doesn't exist
	if err = os.MkdirAll(certDir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create certificate template
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %w", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.AddDate(1, 0, 0) // Valid for 1 year

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Litetable Local Development"},
			CommonName:   "localhost",
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		DNSNames:              []string{"localhost"},
	}

	// Create self-signed certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	// Write certificate to file
	certPath := filepath.Join(certDir, serverCertName)
	certFile, err := os.Create(certPath)
	if err != nil {
		return fmt.Errorf("failed to create certificate file: %w", err)
	}
	defer certFile.Close()

	if err := pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return fmt.Errorf("failed to write certificate: %w", err)
	}

	// Write private key to file
	keyPath := filepath.Join(certDir, serverKeyName)
	keyFile, err := os.Create(keyPath)
	if err != nil {
		return fmt.Errorf("failed to create key file: %w", err)
	}
	defer keyFile.Close()

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	if err := pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyBytes}); err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}

	fmt.Printf("Successfully generated TLS credentials at %s\n", certDir)
	fmt.Printf("Certificate: %s\n", certPath)
	fmt.Printf("Private key: %s\n", keyPath)

	return nil
}
