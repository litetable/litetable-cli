package cmd

import (
	"fmt"
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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
	}
}
