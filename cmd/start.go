package cmd

import (
	"github.com/mritd/pd/helper"
	"github.com/spf13/cobra"
	"os"
)

var startCmd = &cobra.Command{
	Use:     "start",
	Short:   "Start VM",
	Aliases: []string{"run", "open"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			os.Exit(1)
		}
		helper.StartVM(args)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
