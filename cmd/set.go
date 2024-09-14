package cmd

import (
	"github.com/oxio/kv/internal/kv"
	"github.com/oxio/kv/internal/parser"
	"github.com/spf13/cobra"
)

func newSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <file> <key> <value>",
		Short: "Sets a value to a key-value file",
		RunE: func(cmd *cobra.Command, args []string) error {

			file := args[0]
			key := args[1]
			value := args[2]

			repo := kv.NewKvRepo(file)
			item, err := parser.NewItem(key, value)
			if err != nil {
				return err
			}

			return repo.Set(item)
		},
	}
}

func init() {
	rootCmd.AddCommand(newSetCmd())
}
