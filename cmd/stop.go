package cmd

import (
	"github.com/mritd/pd/helper"
	"github.com/spf13/cobra"
	"os"
)

var forceStop bool
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop VMs",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			os.Exit(1)
		}
		helper.StopVM(args, forceStop)
	},
}

func init() {
	stopCmd.PersistentFlags().BoolVarP(&forceStop, "force", "f", false, "Force stop VM")
	rootCmd.AddCommand(stopCmd)
}
