package cmd

import (
	"forklift/api"
	"github.com/spf13/cobra"
)

var forkliftApiCmd = &cobra.Command{
	Use:   "forklift-api",
	Short: "forklift api",
	Long:  `forklift api`,
	RunE:  api.Init,
}

func init() {
	RootCmd.AddCommand(forkliftApiCmd)
}
