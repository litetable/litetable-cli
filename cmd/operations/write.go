package operations

import (
	"context"
	"fmt"
	"github.com/litetable/litetable-cli/cmd/service"
	"github.com/litetable/litetable-cli/internal/server"
	"github.com/spf13/cobra"
	"net/url"
	"time"
)

var (
	writeKey    string
	writeFamily string
	writeQuals  []string
	writeValues []string
	writeTTL    int64

	WriteCmd = &cobra.Command{
		Use:   "write",
		Short: "Write data to the Litetable server",
		Long:  "Write allows you to send data to the Litetable server",
		Example: "litetable write --key=rowKey --family=familyName --qualifier=qual1 --value=val1" +
			" --qualifier=qual2 --value=val2 --ttl=60",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Validate inputs
			if writeKey == "" {
				return fmt.Errorf("key is required")
			}
			if writeFamily == "" {
				return fmt.Errorf("family is required")
			}
			if len(writeQuals) != len(writeValues) {
				return fmt.Errorf("number of qualifiers (%d) must match number of values (%d)",
					len(writeQuals), len(writeValues))
			}
			if len(writeQuals) == 0 {
				return fmt.Errorf("at least one qualifier/value pair is required")
			}
			if writeTTL < 0 {
				return fmt.Errorf("TTL must be a non-negative value")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := writeData(); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		},
	}
)

func init() {
	WriteCmd.Flags().StringVarP(&writeKey, "key", "k", "", "Row key for the write operation")
	WriteCmd.Flags().StringVarP(&writeFamily, "family", "f", "",
		"Column family for the write operation")
	WriteCmd.Flags().StringArrayVarP(&writeQuals, "qualifier", "q", []string{},
		"Qualifiers to read (can be specified multiple times)")
	WriteCmd.Flags().StringArrayVarP(&writeValues, "value", "v", []string{},
		"Values to write (can be specified multiple times, use quotes for values with spaces)")
	WriteCmd.Flags().Int64VarP(&writeTTL, "ttl", "t", 0,
		"Time to live in seconds (0 means no expiration)")
}

func writeData() error {
	conn, err := service.Dial()
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}

	defer conn.Close()

	var quals []server.Qualifier
	// Create the WRITE command with all the qualifier/value pairs
	for i := 0; i < len(writeQuals); i++ {
		// URL encode the value to properly handle spaces and special characters
		encodedValue := url.QueryEscape(writeValues[i])
		quals = append(quals, server.Qualifier{
			Name:  writeQuals[i],
			Value: encodedValue,
		})
	}

	now := time.Now()

	client, err := server.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create server client: %w", err)
	}

	defer client.Close()

	// TODO: fix the ttl stuff after server
	opts := server.WriteParams{
		Key:        writeKey,
		Family:     writeFamily,
		Qualifiers: quals,
	}
	data, err := client.Write(context.Background(), &opts)
	if err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	first := true
	for key, row := range data {
		if !first {
			fmt.Println("--------------------")
			fmt.Println()
		} else {
			fmt.Println()
		}
		first = false
		fmt.Printf("Key: %s\n%s\n", key, row.PrettyPrint())
	}
	fmt.Printf("Row results: %d\n", len(data))
	fmt.Printf("Query duration: %s\n", time.Since(now))

	return nil
}
