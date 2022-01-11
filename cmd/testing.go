package cmd

import (
	"forklift/tests"
	"github.com/spf13/cobra"
)

var testingCmd = &cobra.Command{
	Use:   "testing",
	Short: "test",
	Long:  `test`,
	RunE:  tests.Init,
}

func init() {
	RootCmd.AddCommand(testingCmd)
}
