package cmd

import (
	"github.com/spf13/cobra"
)

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Like docker-machine, But need to create a virtual machine manually",
}

func init() {
	rootCmd.AddCommand(dockerCmd)
}
