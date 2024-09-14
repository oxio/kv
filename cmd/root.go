package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kv",
	Short: "Simple Key-Value storage tool",
}

func Execute() {
	_ = rootCmd.Execute()
}
