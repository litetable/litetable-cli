package cmd

import (
	"encoding/json"
)

// isValidJSON checks if the buffer contains a complete, valid JSON object
func isValidJSON(data []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(data, &js) == nil
}
