package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/litetable/litetable-cli/internal/dir"
	"github.com/spf13/cobra"
)

var stopCommand = &cobra.Command{
	Use:   "stop",
	Short: "Stop the running LiteTable server",
	Long:  "Stop the LiteTable server by sending a SIGTERM signal for graceful shutdown",
	Run: func(cmd *cobra.Command, args []string) {
		if err := stopLiteTable(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCommand)
}

func stopLiteTable() error {
	fmt.Println("⏹️ Stopping LiteTable server...")

	// Get LiteTable directory
	liteTableDir, err := dir.GetLitetableDir()
	if err != nil {
		return fmt.Errorf("failed to get LiteTable directory: %w", err)
	}

	// Check for PID file
	pidFile := filepath.Join(liteTableDir, "litetable.pid")
	if _, err := os.Stat(pidFile); os.IsNotExist(err) {
		return fmt.Errorf("no running LiteTable server found")
	}

	// Read PID from file
	pidBytes, err := os.ReadFile(pidFile)
	if err != nil {
		return fmt.Errorf("failed to read PID file: %w", err)
	}

	pidStr := strings.TrimSpace(string(pidBytes))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return fmt.Errorf("invalid PID in file: %w", err)
	}

	// Find process and send SIGTERM
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process with ID %d: %w", pid, err)
	}

	fmt.Printf("📩 Sending graceful shutdown signal to process %d...\n", pid)

	// Send SIGTERM for graceful shutdown
	if err := process.Signal(syscall.SIGTERM); err != nil {
		// If sending SIGTERM fails, try to remove the PID file anyway
		_ = os.Remove(pidFile)
		return fmt.Errorf("failed to send shutdown signal: %w", err)
	}

	// Wait for a short period to see if a process terminates
	fmt.Println("⏳ Waiting for server to shut down...")

	// Create a timeout channel
	timeout := time.After(5 * time.Second)
	shutdown := make(chan struct{})

	// Check if the process has terminated
	go func() {
		for {
			// Check if process still exists
			if err := process.Signal(syscall.Signal(0)); err != nil {
				// Process has terminated
				close(shutdown)
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Wait for shutdown or timeout
	select {
	case <-shutdown:
		fmt.Println("✅ LiteTable server has been stopped successfully.")
	case <-timeout:
		fmt.Println("⚠️ Timeout waiting for server to stop.")
	}

	// Clean up PID file
	if err := os.Remove(pidFile); err != nil {
		return fmt.Errorf("failed to remove PID file: %w", err)
	}

	return nil
}
