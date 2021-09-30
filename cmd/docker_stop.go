package cmd

import (
	"github.com/mritd/pd/helper"

	"github.com/spf13/cobra"
)

var dockerStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop docker VM",
	Run: func(cmd *cobra.Command, args []string) {
		helper.StopVM([]string{dockerVM},false)
	},
}

func init() {
	dockerStopCmd.PersistentFlags().StringVarP(&dockerVM, "name", "n", "docker", "Docker VM name")
	dockerCmd.AddCommand(dockerStopCmd)
}
