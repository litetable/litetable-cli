package operations

import (
	"context"
	"fmt"
	"github.com/litetable/litetable-cli/internal/server"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

var (
	// Create command options
	families string

	CreateCmd = &cobra.Command{
		Use:     "create",
		Short:   "Create configuration in the Litetable server",
		Long:    "Create configuration elements such as column families in the Litetable server",
		Example: "litetable create -f 'family1, family2, family3'",
		Run: func(cmd *cobra.Command, args []string) {
			createFamilies()
		},
	}
)

func init() {
	// Add flags for create operation
	CreateCmd.Flags().StringVarP(&families, "family", "f", "", "Column families to create (comma-separated)")
	_ = CreateCmd.MarkFlagRequired("family")
}

func createFamilies() {
	start := time.Now()

	var familyParams server.CreateFamilyParams

	// Parse comma-separated families and properly trim whitespace
	familyList := strings.Split(families, ",")
	for _, family := range familyList {
		trimmed := strings.TrimSpace(family)
		// Only add non-empty family names
		if trimmed != "" {
			familyParams.Families = append(familyParams.Families, trimmed)
		}
	}

	// Check if we have any families to create
	if len(familyParams.Families) == 0 {
		fmt.Printf("no valid family names provided")
		return
	}

	// Create a new gRPC client
	client, err := server.NewClient()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	defer func() {
		_ = client.Close()
	}()

	// Create families on the server
	if err = client.CreateFamilies(context.Background(), &familyParams); err != nil {
		fmt.Printf("%w", err)
		return
	}

	fmt.Printf("Created famililes in %s\n", time.Since(start))
}
