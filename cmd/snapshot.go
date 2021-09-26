package cmd

import (
	"github.com/mritd/pd/helper"
	"github.com/spf13/cobra"
	"os"
)

var snapshotCmd = &cobra.Command{
	Use:     "snapshot",
	Short:   "Create VMs Snapshot",
	Aliases: []string{"sp"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			_ = cmd.Help()
			os.Exit(1)
		}
		helper.CreateSnapshot(args[:len(args)-1], args[len(args)-1])
	},
}

func init() {
	rootCmd.AddCommand(snapshotCmd)
}
