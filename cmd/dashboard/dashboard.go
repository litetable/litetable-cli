package dashboard

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/exec"
	"runtime"
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
	// Start a simple HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<html><head><title>LiteTable Dashboard</title></head>"+
			"<body style='font-family: Arial, sans-serif; text-align: center; margin-top: 100px;'>"+
			"<h1>Hello World</h1>"+
			"<p>Welcome to the LiteTable Dashboard</p>"+
			"</body></html>")
	})

	// Start the server in a goroutine
	go func() {
		fmt.Println("Starting dashboard server at http://127.0.0.1:8080")
		if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
			fmt.Printf("Failed to start dashboard: %v\n", err)
			os.Exit(1)
		}
	}()

	// Open browser
	url := "http://127.0.0.1:8080"
	fmt.Printf("Opening browser at %s\n", url)
	openBrowser(url)

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
