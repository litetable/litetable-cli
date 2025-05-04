package dashboard

import (
	"embed"
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

//go:embed web/*
var webContent embed.FS

var (
	dashboardURL = "http:127.0.0.1:8080"
	Command      = &cobra.Command{
		Use:   "dashboard",
		Short: "Open LiteTable dashboard in a browser",
		Long:  "Opens a browser window with the LiteTable dashboard interface",
		Run: func(cmd *cobra.Command, args []string) {
			startDashboard()
		},
	}
)

func startDashboard() {
	// Get web content from embedded files
	webFS, err := fs.Sub(webContent, "web")
	if err != nil {
		fmt.Printf("Failed to load web content: %v\n", err)
		os.Exit(1)
	}

	// Serve static files
	http.Handle("/", http.FileServer(http.FS(webFS)))

	// Start the server in a goroutine
	go func() {
		fmt.Println("Starting dashboard server at http://127.0.0.1:8080")
		if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
			fmt.Printf("Failed to start dashboard: %v\n", err)
			os.Exit(1)
		}
	}()

	// Open browser
	fmt.Printf("Opening browser at %s\n", dashboardURL)
	openBrowser(dashboardURL)

	// Keep the server running
	select {}
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		fmt.Printf("Unsupported platform. Please open browser manually at: %s\n", url)
		return
	}

	if err != nil {
		fmt.Printf("Failed to open browser: %v\n", err)
		fmt.Printf("Please open your browser manually at: %s\n", url)
	}
}
