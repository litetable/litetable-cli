package cmd

import (
	"github.com/spf13/cobra"
)

var (
	applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "apply a mutation to the litetable server",
		Long:  "apply works to mutate the litetable server with writes or deletes",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
)

func init() {
	rootCmd.AddCommand(applyCmd)
}
