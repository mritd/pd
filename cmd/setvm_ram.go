package cmd

import (
	"github.com/mritd/pd/helper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var setVMRAMCmd = &cobra.Command{
	Use:     "ram",
	Short:   "Set VM RAM",
	Aliases: []string{"mem","m"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			_ = cmd.Help()
			os.Exit(1)
		}
		size, err := strconv.Atoi(args[len(args)-1])
		if err != nil {
			logrus.Fatalf("Invalid number of RAM: %s", args[len(args)-1])
		}
		helper.SetVMRAM(args[:len(args)-1], size)
	},
}

func init() {
	setVMCmd.AddCommand(setVMRAMCmd)
}
