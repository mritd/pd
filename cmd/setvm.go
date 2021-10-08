package cmd

import (
	"github.com/spf13/cobra"
)

var setVMCmd = &cobra.Command{
	Use:     "setvm",
	Short:   "Set VM Config",
	Aliases: []string{"set"},
}

func init() {
	rootCmd.AddCommand(setVMCmd)
}
