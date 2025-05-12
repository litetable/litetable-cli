package operations

import (
	"context"
	"errors"
	"fmt"
	"github.com/litetable/litetable-cli/internal/server"
	"github.com/spf13/cobra"
	"time"
)

var (
	// Read command options
	readKey       string
	readKeyPrefix string
	readRegex     string
	readFamily    string
	readQualifier []string
	readLatest    int

	ReadCmd = &cobra.Command{
		Use:   "read",
		Short: "Read data from the Litetable server",
		Long:  "Read allows you to retrieve data from the Litetable server",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Count how many key selectors are provided
			selectors := 0
			if readKey != "" {
				selectors++
			}
			if readKeyPrefix != "" {
				selectors++
			}
			if readRegex != "" {
				selectors++
			}

			if selectors != 1 {
				return fmt.Errorf("exactly one of --key (-k), --keyPrefix (-p), or --regex (-r) must be provided")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			readData()
		},
	}
)

func init() {
	// Add flags for read operation
	ReadCmd.Flags().StringVarP(&readKey, "key", "k", "", "Row key to read")
	ReadCmd.Flags().StringVarP(&readKeyPrefix, "keyPrefix", "p", "",
		"Read all row-keys with this prefix")
	ReadCmd.Flags().StringVarP(&readRegex, "regex", "r", "",
		"Read all row-keys matching this regex pattern")
	ReadCmd.Flags().StringVarP(&readFamily, "family", "f", "", "Column family to read")
	ReadCmd.Flags().StringArrayVarP(&readQualifier, "qualifier", "q", []string{}, "Qualifiers to read (can be specified multiple times)")
	ReadCmd.Flags().IntVarP(&readLatest, "latest", "l", 0, "Number of latest versions to return")

	// Mark required flags - removing the required mark for key
	_ = ReadCmd.MarkFlagRequired("family")
}

func readData() {
	now := time.Now()

	// Build the READ command based on which selector is provided
	var qualifiers []string
	var queryType server.QueryType
	var queryKey string
	if readKey != "" {
		queryType = server.Read
		queryKey = readKey
	} else if readKeyPrefix != "" {
		queryType = server.ReadPrefix
		queryKey = readKeyPrefix
	} else if readRegex != "" {
		// Create a properly formatted regex pattern that escapes special characters
		// and wraps the user input with ".*" for substring matching
		queryType = server.ReadRegex
		queryKey = fmt.Sprintf(".*%s.*", readRegex)
	}

	for _, q := range readQualifier {
		qualifiers = append(qualifiers, q)
	}

	client, err := server.NewClient()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	defer func(client *server.GrpcClient) {
		_ = client.Close()
	}(client)

	opts := server.ReadParams{
		Key:        queryKey,
		QueryType:  queryType,
		Family:     readFamily,
		Qualifiers: qualifiers,
		Latest:     int32(readLatest),
	}
	data, err := client.Read(context.Background(), &opts)
	if err != nil {
		if errors.Is(err, server.ErrRowNotFound) {
			fmt.Println("row not found")
			return
		}
		fmt.Printf("failed to read data: %w", err)
		return
	}

	first := true
	// Print the rows
	for _, row := range data {
		if !first {
			fmt.Println("--------------------")
			fmt.Println()
		}
		first = false
		fmt.Printf("%s\n", row.PrettyPrint())
	}
	fmt.Printf("Row results: %d\n", len(data))
	fmt.Printf("Query duration: %s\n", time.Since(now))
}
