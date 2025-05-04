package cmd

import (
	"fmt"
	"github.com/litetable/litetable-cli/internal/litetable"
	"github.com/spf13/cobra"
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Display LiteTable version information",
	Long:  "Show the installed version of LiteTable server from configuration",
	Run: func(cmd *cobra.Command, args []string) {
		version, err := litetable.GetFromConfig(litetable.ServerVersionKey)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("LiteTable Server Version: %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCommand)
}
