package cmd

import (
	"github.com/spf13/cobra"
)

var snapshotCmd = &cobra.Command{
	Use:     "snapshot",
	Short:   "Snapshot Manager",
	Aliases: []string{"sp"},
}

func init() {
	rootCmd.AddCommand(snapshotCmd)
}
