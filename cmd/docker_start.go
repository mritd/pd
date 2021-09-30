package cmd

import (
	"github.com/mritd/pd/helper"

	"github.com/spf13/cobra"
)

var dockerVM string
var dockerBindingHome bool
var dockerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start docker VM",
	Run: func(cmd *cobra.Command, args []string) {
		helper.StartDockerVM(dockerVM, dockerBindingHome)
	},
}

func init() {
	dockerStartCmd.PersistentFlags().StringVarP(&dockerVM, "name", "n", "docker", "Docker VM name")
	dockerStartCmd.PersistentFlags().BoolVar(&dockerBindingHome, "bind-home", false, "Binding User Home to VM `/root` dir")
	dockerCmd.AddCommand(dockerStartCmd)
}
