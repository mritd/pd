package cmd

import (
	"github.com/mritd/pd/helper"
	"github.com/spf13/cobra"
)

var listAll bool
var lsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List VMs",
	Run: func(cmd *cobra.Command, args []string) {
		helper.ListVMs(listAll)
	},
}

func init() {
	lsCmd.PersistentFlags().BoolVarP(&listAll, "all", "a", false, "List all vms")
	rootCmd.AddCommand(lsCmd)
}
