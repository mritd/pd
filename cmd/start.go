package cmd

import (
	"github.com/mritd/pd/helper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var fake bool

var startCmd = &cobra.Command{
	Use:     "start",
	Short:   "Start VM",
	Aliases: []string{"run", "open"},
	PreRun: func(cmd *cobra.Command, args []string) {
		if fake {
			if err := helper.FakeDate(); err != nil {
				logrus.Fatal("Failed to update system date")
			}
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		if fake {
			if err := helper.RestoreDate(); err != nil {
				logrus.Fatal("Failed to restore system date")
			}
		}
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		vms, _ := helper.ListVMInfo(true)
		var ss []string
		for _, vm := range vms {
			if vm.Status == "stopped" {
				ss = append(ss, vm.Name+"\t"+vm.UUID)
			}
		}
		return ss, cobra.ShellCompDirectiveDefault
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			os.Exit(1)
		}
		helper.StartVM(args)
	},
}

func init() {
	startCmd.PersistentFlags().BoolVar(&fake, "fake", false, "Start VM With Fake System Date")
	rootCmd.AddCommand(startCmd)
}
