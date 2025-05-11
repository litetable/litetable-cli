package operations

import (
	"context"
	"fmt"
	"github.com/litetable/litetable-cli/internal/server"
	"github.com/spf13/cobra"
	"time"
)

var (
	// Delete command options
	deleteKey       string
	deleteFamily    string
	deleteQualifier []string
	deleteTTL       int64
	deleteFrom      int64

	DeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete data from the Litetable server",
		Long:  "Delete allows you to remove data from the Litetable server",
		Run: func(cmd *cobra.Command, args []string) {
			// Execute the delete operation
			if err := deleteData(); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		},
	}
)

func init() {
	// Add flags for delete operation
	DeleteCmd.Flags().StringVarP(&deleteKey, "key", "k", "", "Row key to delete (required)")
	DeleteCmd.Flags().StringVarP(&deleteFamily, "family", "f", "", "Column family to delete")
	DeleteCmd.Flags().StringArrayVarP(&deleteQualifier, "qualifier", "q", []string{}, "Qualifiers to delete (can be specified multiple times)")
	DeleteCmd.Flags().Int64Var(&deleteTTL, "ttl", 0, "Time-to-live in seconds for tombstone entries")
	DeleteCmd.Flags().Int64Var(&deleteFrom, "from", 0, "Starting position for deletion in the map")

	// Mark required flags
	_ = DeleteCmd.MarkFlagRequired("key")
}

func deleteData() error {
	now := time.Now()

	var qualifiers []string

	for _, q := range deleteQualifier {
		qualifiers = append(qualifiers, q)
	}

	opts := &server.DeleteParams{
		Key:        deleteKey,
		Family:     deleteFamily,
		Qualifiers: deleteQualifier,
	}
	if deleteFrom > 0 {
		opts.From = deleteFrom
	}
	if deleteTTL > 0 {
		opts.TTL = int32(deleteTTL)
	}

	client, err := server.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create server client: %w", err)
	}

	defer func(client *server.GrpcClient) {
		_ = client.Close()
	}(client)

	if err = client.Delete(context.Background(), opts); err != nil {
		return fmt.Errorf("failed to delete data: %w", err)
	}

	fmt.Println("Delete successful in", time.Since(now))
	return nil
}
