package litetable

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

const (
	CLIVersion  = "0.1.0"
	DatabaseURL = "https://github.com/litetable/litetable-db"
)

// TimestampedValue stores a value with its timestamp
type TimestampedValue struct {
	Value     []byte    `json:"-"`         // Internal binary representation
	RawValue  string    `json:"value"`     // Base64 encoded value from JSON
	Timestamp time.Time `json:"-"`         // Parsed timestamp
	RawTime   string    `json:"timestamp"` // String timestamp from JSON
}

// UnmarshalJSON implements custom unmarshalling for TimestampedValue
func (tv *TimestampedValue) UnmarshalJSON(data []byte) error {
	// Temporary struct for parsing
	temp := struct {
		Value     string `json:"value"`
		Timestamp string `json:"timestamp"`
	}{}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Store raw values
	tv.RawValue = temp.Value
	tv.RawTime = temp.Timestamp

	// Decode base64 value
	decoded, err := base64.StdEncoding.DecodeString(temp.Value)
	if err != nil {
		return fmt.Errorf("invalid base64 value: %w", err)
	}
	tv.Value = decoded

	// Parse timestamp
	timestamp, err := time.Parse(time.RFC3339Nano, temp.Timestamp)
	if err != nil {
		return fmt.Errorf("invalid timestamp format: %w", err)
	}
	tv.Timestamp = timestamp

	return nil
}

// GetString returns the decoded value as a string
func (tv *TimestampedValue) GetString() string {
	return string(tv.Value)
}

// VersionedQualifier maps qualifiers to their timestamped values
type VersionedQualifier map[string][]TimestampedValue

// Row represents a row of data
type Row struct {
	Key     string                        `json:"key"`
	Columns map[string]VersionedQualifier `json:"cols"` // family → qualifier → []TimestampedValue
}

// PrettyPrint formats the row in a human-readable way
func (r *Row) PrettyPrint() string {
	var result string
	result += fmt.Sprintf("rowKey: %s\n", r.Key)

	for family, qualifiers := range r.Columns {
		result += fmt.Sprintf("family: %s\n", family)

		for qualifier, values := range qualifiers {
			result += fmt.Sprintf("  qualifier: %s\n", qualifier)

			for i, v := range values {
				result += fmt.Sprintf("    value %d: %s (timestamp: %s)\n",
					i+1, v.GetString(), v.RawTime)
			}
		}
	}

	return result
}
