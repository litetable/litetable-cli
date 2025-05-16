package dashboard

import (
	"embed"
	"fmt"
	"github.com/litetable/litetable-cli/internal/server"
	"github.com/spf13/cobra"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

//go:embed web/*
var webContent embed.FS

const (
	dashboardHost = "127.0.0.1"
	dashboardPort = "7654"
)

var (
	Command = &cobra.Command{
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

	litetableClient, err := server.NewClient()
	if err != nil {
		fmt.Printf("Failed to create LiteTable client: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		_ = litetableClient.Close()
	}()

	ltHandler := &handler{
		server: litetableClient,
	}

	// Serve static files
	http.Handle("/", http.FileServer(http.FS(webFS)))

	// Serve handlers
	http.Handle("POST /query", http.HandlerFunc(ltHandler.query))
	http.Handle("GET /families", http.HandlerFunc(ltHandler.getFamilies))

	addr := fmt.Sprintf("%s:%s", dashboardHost, dashboardPort)

	// Start the server in a goroutine
	go func() {
		// Start the HTTP server
		fmt.Printf("Starting dashboard server at http://%s\n", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			fmt.Printf("Failed to start dashboard: %v\n", err)
			os.Exit(1)
		}
	}()

	// Open browser
	fmt.Printf("Opening browser at %s\n", addr)
	openBrowser(fmt.Sprintf("http://%s", addr))

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
