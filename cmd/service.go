package cmd

import (
	"github.com/litetable/litetable-cli/cmd/service"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage LiteTable service",
	Long:  "Commands for managing the LiteTable server",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	// Add the service command to the root command
	rootCmd.AddCommand(serviceCmd)

	// Add subcommands to the service command
	serviceCmd.AddCommand(service.StartCommand)
	serviceCmd.AddCommand(service.StopCommand)
	// serviceCmd.AddCommand(service.CredentialsCmd)
	serviceCmd.AddCommand(service.HealthCmd)
}
