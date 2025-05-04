package cmd

import (
	"fmt"
	"github.com/litetable/litetable-cli/cmd/dashboard"
	"github.com/litetable/litetable-cli/cmd/operations"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "litetable",
		Example: "litetable --help\n\nlitetable service init",
		Short:   "A CLI tool for interacting with litetable",
		Long: "Litetable is a high-performance key-value store designed with local" +
			" development in mind. Proudly written in pure Go.\n",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
)

func init() {
	// Add operations commands to the root command
	rootCmd.AddCommand(operations.CreateCmd)
	rootCmd.AddCommand(operations.ReadCmd)
	rootCmd.AddCommand(operations.WriteCmd)
	rootCmd.AddCommand(operations.DeleteCmd)
	rootCmd.AddCommand(dashboard.DashboardCommand)

	rootCmd.AddCommand(serviceCmd)
	rootCmd.AddCommand(UpdateCmd)

	rootCmd.AddCommand(uninstallCommand)
	rootCmd.AddCommand(versionCommand)

	rootCmd.AddCommand(wipeCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
	}
}
