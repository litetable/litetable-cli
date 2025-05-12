package service

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"github.com/litetable/litetable-cli/internal/dir"
	"github.com/spf13/cobra"
)

var StartCommand = &cobra.Command{
	Use:   "start",
	Short: "Start the LiteTable server",
	Long:  "Start the LiteTable server if installed, otherwise prompt to run init",
	Run: func(cmd *cobra.Command, args []string) {
		if err := startLiteTable(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	},
}

func startLiteTable() error {
	fmt.Println("🚀 Starting LiteTable server...")

	// Check if the server is already running
	isRunning, pid, err := checkProcessRunning()
	if err != nil {
		fmt.Printf("⚠️  Warning: Could not determine if server is running: %v\n", err)
	} else if isRunning {
		fmt.Println("✅  LiteTable server is already running.")
		return nil
	}

	// Get LiteTable directory
	liteTableDir, err := dir.GetLitetableDir()
	if err != nil {
		return fmt.Errorf("failed to get LiteTable directory: %w", err)
	}

	// Check if LiteTable is installed
	binPath := filepath.Join(liteTableDir, "bin", serverBin)
	if runtime.GOOS == "windows" {
		binPath += ".exe"
	}

	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		fmt.Println("\n⚠️ LiteTable server is not installed.")
		fmt.Print("Would you like to run the init command now? (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return initLiteTable()
		}
		return fmt.Errorf("server not installed. Run 'litetable-cli init' to install")
	}

	// Read the config file to get server binary location
	configPath := filepath.Join(liteTableDir, "litetable.conf")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// If config doesn't exist, use the default binary path
		fmt.Println("⚠️ Configuration file not found, using default binary path")
	} else {
		// Read the config file to get the server binary path
		configBytes, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("failed to read config file: %w", err)
		}

		// Parse the config file to find server_binary
		configLines := strings.Split(string(configBytes), "\n")
		for _, line := range configLines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "server_binary") {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					binPath = strings.TrimSpace(parts[1])
				}
				break
			}
		}
	}

	// Create a log file
	logFile := filepath.Join(liteTableDir, "litetable.log")
	logFileHandle, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	defer logFileHandle.Close()

	// Start the server
	fmt.Printf("📡 Running LiteTable server from: %s\n", binPath)
	fmt.Printf("📝 Logs will be written to: %s\n", logFile)

	// Create a new command to start the server
	serverCmd := exec.Command(binPath)
	serverCmd.Stdout = logFileHandle
	serverCmd.Stderr = logFileHandle

	// Detach the process from the terminal
	serverCmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if err := serverCmd.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	pid = serverCmd.Process.Pid
	pidFile := filepath.Join(liteTableDir, "litetable.pid")
	if err := os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("failed to write PID file: %w", err)
	}

	fmt.Printf("✅  LiteTable server started with PID: %d\n", pid)
	return nil
}

func checkProcessRunning() (bool, int, error) {
	// Get LiteTable directory
	liteTableDir, err := dir.GetLitetableDir()
	if err != nil {
		return false, 0, fmt.Errorf("failed to get LiteTable directory: %w", err)
	}

	// Check if PID file exists
	pidFile := filepath.Join(liteTableDir, "litetable.pid")
	pidBytes, err := os.ReadFile(pidFile)
	if err != nil {
		if os.IsNotExist(err) {
			return false, 0, nil // PID file doesn't exist, so server isn't running
		}
		return false, 0, fmt.Errorf("failed to read PID file: %w", err)
	}

	// Parse PID
	pid, err := strconv.Atoi(strings.TrimSpace(string(pidBytes)))
	if err != nil {
		return false, 0, fmt.Errorf("failed to parse PID: %w", err)
	}

	// Check if process is running
	process, err := os.FindProcess(pid)
	if err != nil {
		return false, pid, fmt.Errorf("failed to find process: %w", err)
	}

	// On Unix systems, FindProcess always succeeds, so we need to send signal 0 to check
	// if the process exists. On Windows, FindProcess only succeeds if the process exists.
	if err := process.Signal(syscall.Signal(0)); err != nil {
		return false, pid, nil // Process doesn't exist
	}

	return true, pid, nil // Process exists
}
