package litetable

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	DatabaseURL   = "https://github.com/litetable/litetable-db"
	CLIURL        = "https://github.com/litetable/litetable-cli"
	CLIInstallURL = "https://raw.githubusercontent.com/litetable/litetable-cli/main/install.sh"
)

// TimestampedValue stores a value with its timestamp
type TimestampedValue struct {
	Value     []byte `json:"value"`     // Internal binary representation
	RawValue  string `json:"-"`         // Base64 encoded value from JSON
	Timestamp int64  `json:"timestamp"` // Parsed timestamp
}

// UnmarshalJSON implements custom unmarshalling for TimestampedValue
func (tv *TimestampedValue) UnmarshalJSON(data []byte) error {
	// Temporary struct for parsing
	temp := struct {
		Value     string `json:"value"`
		Timestamp int64  `json:"timestamp_unix"`
	}{
		Value:     tv.RawValue,
		Timestamp: tv.Timestamp,
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Store raw values
	tv.RawValue = temp.Value
	tv.Timestamp = temp.Timestamp

	// Decode base64 value
	decoded, err := base64.StdEncoding.DecodeString(temp.Value)
	if err != nil {
		return fmt.Errorf("invalid base64 value: %w", err)
	}
	tv.Value = decoded

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
				// Try to URL-decode the value
				rawValue := v.GetString()
				decodedValue, err := url.QueryUnescape(rawValue)
				if err != nil {
					// If it fails to decode, use the original value
					decodedValue = rawValue
				}

				result += fmt.Sprintf("    value %d: %s, timestamp: %d\n",
					i+1, decodedValue, v.Timestamp)
			}
		}
	}

	return result
}
