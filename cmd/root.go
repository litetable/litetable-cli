package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "litetable",
		Example: "litetable --help",
		Short:   "A CLI tool for interacting with litetable",
		Long:    "Litetable is a high-performance key-value store designed for local development. Proudly written in pure Go.\n",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
