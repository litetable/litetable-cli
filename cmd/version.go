package cmd

import (
	"fmt"
	"github.com/litetable/litetable-cli/internal/litetable"
	"github.com/spf13/cobra"
)

var (
	versionCommand = &cobra.Command{
		Use:   "version",
		Short: "Display LiteTable version information",
		Long:  "Show the installed version of LiteTable server from configuration",
		Run: func(cmd *cobra.Command, args []string) {
			version := "not set"
			foundVersion, err := litetable.GetFromConfig(litetable.ServerVersionKey)
			if err != nil {
				fmt.Println("\nRun \033[0;33m`litetable service init`\033[0m to configure the server.")
			} else {
				version = foundVersion
			}

			fmt.Printf("\nLiteTable CLI Version: %s\n", litetable.CLIVersion)
			fmt.Printf("LiteTable Server Version: %s\n", version)
		},
	}
)
