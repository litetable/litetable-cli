package dashboard

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/litetable/litetable-cli/cmd/service"
)

func queryHandlerFunc() http.HandlerFunc {
	return queryHandler
}

type payload struct {
	Query string `json:"query"`
}

// queryHandler only handles the POST method to pass the query to the UI. All the parsing
// is done in the UI.
func queryHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Only POST method is supported",
		})
		return
	}

	// Decode the JSON payload
	var p payload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to decode JSON payload: %v", err),
		})
		return
	}

	// Check if query is empty
	if p.Query == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Query string cannot be empty",
		})
		return
	}

	// Connect to the litetable server
	conn, err := service.Dial()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to dial server: %v", err),
		})
		return
	}
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	// Set a reasonable timeout
	timeout := time.Now().Add(10 * time.Second)
	if err := conn.SetReadDeadline(timeout); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to set read deadline: %v", err),
		})
		return
	}

	// Send the query to the server
	if _, err = conn.Write([]byte(p.Query)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to send query to server: %v", err),
		})
		return
	}

	// Read the response
	var fullResponse []byte
	buffer := make([]byte, 4096)

	for {
		n, err := conn.Read(buffer)
		if n > 0 {
			fullResponse = append(fullResponse, buffer[:n]...)

			// Check if we have a complete JSON object
			if len(fullResponse) > 0 && isValidJSON(fullResponse) {
				break
			}
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Error reading response: %v", err),
			})
			return
		}

		// Extend deadline for each successful read
		if err := conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Failed to extend read deadline: %v", err),
			})
			return
		}
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Write the raw server response directly to the client
	_, _ = w.Write(fullResponse)
}

// isValidJSON checks if the buffer contains a complete, valid JSON object
func isValidJSON(data []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(data, &js) == nil
}
