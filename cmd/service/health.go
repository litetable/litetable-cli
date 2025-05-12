package service

import (
	"context"
	"fmt"
	"github.com/litetable/litetable-cli/internal/litetable"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"time"
)

// HealthCmd represents the health command
var HealthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check if the LiteTable server is running",
	Long:  `Sends a PING request to the LiteTable server to verify it's operational.`,
	Run: func(cmd *cobra.Command, args []string) {
		checkServerHealth()
	},
}

func checkServerHealth() {
	fmt.Println("🔍 Checking LiteTable server health...")
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	serverAddress, _ := litetable.GetFromConfig(litetable.ServerAddress)
	port, _ := litetable.GetFromConfig(litetable.ServerPort)

	// Build the URL
	url := fmt.Sprintf("http://%s:%s/health", serverAddress, port)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("failed to create request: %v", err)
		return
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("server health check failed: %v", err)
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to read response body: %v`", err)
		return
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("server health check returned non-OK status: %d - %s",
			resp.StatusCode, string(body))
		return
	}

	fmt.Println("✅  LiteTable server is healthy!")
	fmt.Printf("%s\n", string(body))

}
