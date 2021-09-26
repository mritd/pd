package cmd

import (
	"github.com/mritd/pd/helper"
	"os"

	"github.com/spf13/cobra"
)

var switchSnapshotCmd = &cobra.Command{
	Use:   "switchSnapshot",
	Short: "Switch VM Snapshot",
	Aliases: []string{"ssp"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			_ = cmd.Help()
			os.Exit(1)
		}
		helper.SwitchSnapshot(args[:len(args)-1], args[len(args)-1])
	},
}

func init() {
	rootCmd.AddCommand(switchSnapshotCmd)
}
