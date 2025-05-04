package cmd

import (
	"github.com/litetable/litetable-cli/cmd/operations"
)

func init() {
	// Add operations commands to the root command
	rootCmd.AddCommand(operations.CreateCmd)
	rootCmd.AddCommand(operations.ReadCmd)
	rootCmd.AddCommand(operations.WriteCmd)
	rootCmd.AddCommand(operations.DeleteCmd)
}
