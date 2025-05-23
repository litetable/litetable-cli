package cmd

import (
	"github.com/litetable/litetable-cli/cmd/service"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage the LiteTable service",
	Long:  "Commands for managing the LiteTable server",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	// Add subcommands to the service command
	serviceCmd.AddCommand(service.InitCommand)
	serviceCmd.AddCommand(service.StartCommand)
	serviceCmd.AddCommand(service.StopCommand)
	// serviceCmd.AddCommand(service.CredentialsCmd)
	serviceCmd.AddCommand(service.UpdateCommand)
	serviceCmd.AddCommand(service.HealthCmd)
}
