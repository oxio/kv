package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "kvf",
	Short:   "Simple Key-Value storage tool",
	Version: "1.2.2",
}

func Execute() {
	_ = rootCmd.Execute()
}
