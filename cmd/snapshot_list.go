package cmd

import (
	"github.com/mritd/pd/helper"
	"github.com/spf13/cobra"
	"os"
)

var listSnapshotCmd = &cobra.Command{
	Use:     "list",
	Short:   "List Snapshot",
	Aliases: []string{"l", "ls"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			_ = cmd.Help()
			os.Exit(1)
		}
		helper.ListSnapshot(args[0])
	},
}

func init() {
	snapshotCmd.AddCommand(listSnapshotCmd)
}
