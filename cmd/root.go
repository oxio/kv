package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kvf",
	Short: "Simple Key-Value storage tool",
}

func Execute() {
	_ = rootCmd.Execute()
}
