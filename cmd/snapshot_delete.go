package cmd

import (
	"github.com/mritd/pd/helper"
	"github.com/spf13/cobra"
	"os"
)

var deleteSnapshotCmd = &cobra.Command{
	Use:     "deletesp",
	Short:   "Delete VMs Snapshot",
	Aliases: []string{"dsp"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			_ = cmd.Help()
			os.Exit(1)
		}
		helper.DeleteSnapshot(args[:len(args)-1], args[len(args)-1])
	},
}

func init() {
	rootCmd.AddCommand(deleteSnapshotCmd)
}
