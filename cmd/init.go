package cmd

import (
	"bufio"
	"fmt"
	"github.com/litetable/litetable-cli/internal/dir"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	supportedVersion = "v0.1.1"
)

var (
	forceInit   bool
	autostart   bool
	serverRepo  = "https://github.com/litetable/litetable-db"
	serverBin   = "litetable-server"
	initCommand = &cobra.Command{
		Use:   "init",
		Short: "Initialize LiteTable database",
		Long:  "Pull, build, and configure the latest version of LiteTable database server",
		Run: func(cmd *cobra.Command, args []string) {
			if err := initLiteTable(); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(initCommand)
	initCommand.Flags().BoolVarP(&forceInit, "force", "f", false, "Force reinstallation if already installed")
	initCommand.Flags().BoolVarP(&autostart, "autostart", "a", false, "Configure server to start automatically")
}

func initLiteTable() error {
	fmt.Println("🚀 Welcome to LiteTable Setup!")
	fmt.Println("Let's get your database ready to use.")

	// Get LiteTable directory
	liteTableDir, err := dir.GetLitetableDir()
	if err != nil {
		return fmt.Errorf("failed to get LiteTable directory: %w", err)
	}

	// Check if LiteTable is already installed
	binPath := filepath.Join(liteTableDir, "bin", serverBin)
	if runtime.GOOS == "windows" {
		binPath += ".exe"
	}

	// serverExists := false
	if _, err := os.Stat(binPath); err == nil {
		// serverExists = true
		if !forceInit {
			fmt.Println("\n⚠️  LiteTable server appears to be already installed.")
			fmt.Print("Would you like to reinstall? (y/n): ")
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" && response != "yes" {
				fmt.Println("Installation canceled. Your existing setup remains unchanged.")
				return nil
			}
		}
	}

	// Create the necessary directories
	binDir := filepath.Join(liteTableDir, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	// And in your initLiteTable function:
	if !checkGitInstalled() {
		return fmt.Errorf("git is not installed. Please install Git (https://git-scm.com/downloads) and try again")
	}
	fmt.Println("✅ Git installation detected")

	// Check for Go installation
	fmt.Println("\n📋 Checking prerequisites...")
	if !checkGoInstalled() {
		return fmt.Errorf("go is not installed. Please install Go (https://go.dev/doc/install) and try again")
	}
	fmt.Println("✅ Go installation detected")

	// Clone/update repo and build
	fmt.Println("\n📥 Downloading latest LiteTable server...")
	tempDir, err := os.MkdirTemp("", "litetable-build")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
		}
	}(tempDir)

	// Clone the repository
	fmt.Println("\n📥 Cloning LiteTable server repository...")
	gitCloneCmd := exec.Command("git", "clone", "--depth", "1", "--branch", supportedVersion, serverRepo, tempDir)
	gitCloneCmd.Stdout = os.Stdout
	gitCloneCmd.Stderr = os.Stderr
	if err := gitCloneCmd.Run(); err != nil {
		return fmt.Errorf("failed to clone server repository: %w", err)
	}

	// Build the server
	fmt.Println("\n🔨 Building server...")
	buildCmd := exec.Command("go", "build", "-o", binPath)
	buildCmd.Dir = tempDir // Run the build command in the cloned repository directory
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build server: %w", err)
	}

	// Make server executable (especially important for Unix systems)
	if runtime.GOOS != "windows" {
		if err := os.Chmod(binPath, 0755); err != nil {
			return fmt.Errorf("failed to make server executable: %w", err)
		}
	}

	// Generate TLS certificates if they don't exist
	certDir := liteTableDir
	certPath := filepath.Join(certDir, serverCertName)
	// keyPath := filepath.Join(certDir, serverKeyName)

	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		fmt.Println("\n🔐 Generating TLS certificates...")
		if err := generateCredentials(); err != nil {
			return fmt.Errorf("failed to generate credentials: %w", err)
		}
	}

	// Setup autostart if requested
	if autostart {
		fmt.Println("\n⚙️ Setting up autostart...")
		if err := setupAutostart(binPath); err != nil {
			return fmt.Errorf("failed to configure autostart: %w", err)
		}
	}

	// Write a configuration file
	configFile := filepath.Join(liteTableDir, "litetable.conf")
	if err := writeConfigFile(configFile); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	// Success message
	fmt.Println("\n✅ LiteTable setup complete!")
	fmt.Printf("Server installed at: %s\n", binPath)
	if !autostart {
		fmt.Println("\nTo start the server manually:")
		fmt.Printf("  %s\n", binPath)
		fmt.Println("\nOr use the CLI command:")
		fmt.Println("  litetable-cli service start")
	} else {
		fmt.Println("\nServer configured to start automatically")
		fmt.Println("You can manage the service with:")
		fmt.Println("  litetable-cli service status|start|stop|restart")
	}

	return nil
}

func checkGitInstalled() bool {
	cmd := exec.Command("git", "--version")
	return cmd.Run() == nil
}

func checkGoInstalled() bool {
	cmd := exec.Command("go", "version")
	return cmd.Run() == nil
}

func writeConfigFile(path string) error {
	content := `# LiteTable Server Configuration
port = 9090
# Add other configuration options as needed
`
	return os.WriteFile(path, []byte(content), 0644)
}

func setupAutostart(serverPath string) error {
	switch runtime.GOOS {
	case "darwin":
		return setupDarwinAutostart(serverPath)
	case "linux":
		return setupLinuxAutostart(serverPath)
	case "windows":
		return setupWindowsAutostart(serverPath)
	default:
		return fmt.Errorf("autostart not supported on %s", runtime.GOOS)
	}
}

func setupDarwinAutostart(serverPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	launchAgentsDir := filepath.Join(homeDir, "Library", "LaunchAgents")
	if err := os.MkdirAll(launchAgentsDir, 0755); err != nil {
		return err
	}

	plistPath := filepath.Join(launchAgentsDir, "com.litetable.server.plist")
	plistContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>com.litetable.server</string>
	<key>ProgramArguments</key>
	<array>
		<string>%s</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
	<key>KeepAlive</key>
	<true/>
	<key>StandardOutPath</key>
	<string>%s</string>
	<key>StandardErrorPath</key>
	<string>%s</string>
</dict>
</plist>`, serverPath,
		filepath.Join(homeDir, ".litetable", "server.log"),
		filepath.Join(homeDir, ".litetable", "server.err"))

	if err := os.WriteFile(plistPath, []byte(plistContent), 0644); err != nil {
		return err
	}

	loadCmd := exec.Command("launchctl", "load", plistPath)
	return loadCmd.Run()
}

func setupLinuxAutostart(serverPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Simple systemd user service setup
	systemdDir := filepath.Join(homeDir, ".config", "systemd", "user")
	if err := os.MkdirAll(systemdDir, 0755); err != nil {
		return err
	}

	servicePath := filepath.Join(systemdDir, "litetable-server.service")
	serviceContent := fmt.Sprintf(`[Unit]
Description=LiteTable Database Server
After=network.target

[Service]
ExecStart=%s
Restart=on-failure
StandardOutput=append:%s
StandardError=append:%s

[Install]
WantedBy=default.target
`, serverPath,
		filepath.Join(homeDir, ".litetable", "server.log"),
		filepath.Join(homeDir, ".litetable", "server.err"))

	if err := os.WriteFile(servicePath, []byte(serviceContent), 0644); err != nil {
		return err
	}

	// Enable and start the service
	enableCmd := exec.Command("systemctl", "--user", "enable", "litetable-server.service")
	if err := enableCmd.Run(); err != nil {
		return err
	}

	startCmd := exec.Command("systemctl", "--user", "start", "litetable-server.service")
	return startCmd.Run()
}

func setupWindowsAutostart(serverPath string) error {
	// Create a batch file in the startup folder
	startupDir, err := getWindowsStartupDir()
	if err != nil {
		return err
	}

	batchPath := filepath.Join(startupDir, "LiteTableServer.bat")
	logDir, err := dir.GetLitetableDir()
	if err != nil {
		return err
	}

	batchContent := fmt.Sprintf(`@echo off
start "" /B "%s" > "%s" 2>&1
`, serverPath, filepath.Join(logDir, "server.log"))

	return os.WriteFile(batchPath, []byte(batchContent), 0644)
}

func getWindowsStartupDir() (string, error) {
	cmd := exec.Command("powershell", "-Command",
		"[Environment]::GetFolderPath('Startup')")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
