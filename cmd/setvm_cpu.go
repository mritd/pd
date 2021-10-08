package cmd

import (
	"github.com/mritd/pd/helper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var setVMCPUCmd = &cobra.Command{
	Use:     "cpu",
	Short:   "Set VM CPU",
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			_ = cmd.Help()
			os.Exit(1)
		}
		count, err := strconv.Atoi(args[len(args)-1])
		if err != nil {
			logrus.Fatalf("Invalid number of CPUs: %s", args[len(args)-1])
		}
		helper.SetVMCPU(args[:len(args)-1], count)
	},
}

func init() {
	setVMCmd.AddCommand(setVMCPUCmd)
}
