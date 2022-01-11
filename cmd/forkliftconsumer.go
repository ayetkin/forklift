package cmd

import (
	"forklift/internal/consumer"
	"github.com/spf13/cobra"
)

var forkliftConsumerCmd = &cobra.Command{
	Use:   "forklift-consumer",
	Short: "forklift consumer",
	Long:  `forklift consumer`,
	RunE:  consumer.Init,
}

func init() {
	RootCmd.AddCommand(forkliftConsumerCmd)
}
