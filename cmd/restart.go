package cmd

import (
	"github.com/mritd/pd/helper"
	"os"

	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart VM",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			os.Exit(1)
		}
		helper.StopVM(args, forceStop)
		helper.StartVM(args)
	},
}

func init() {
	restartCmd.PersistentFlags().BoolVarP(&forceStop, "force", "f", false, "Force stop VM")
	rootCmd.AddCommand(restartCmd)
}
