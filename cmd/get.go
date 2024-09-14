package cmd

import (
	"github.com/oxio/kv/internal/kv"
	"github.com/oxio/kv/internal/parser"
	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	const (
		defaultValFlag = "default"
	)
	var defaultVal *string

	cmd := &cobra.Command{
		Use:   "get <file> <key> [--default|-d value] ",
		Short: "Gets a value from the key-value file",
		RunE: func(cmd *cobra.Command, args []string) error {
			var item *parser.Item
			var err error
			file := args[0]
			key := args[1]

			repo := kv.NewKvRepo(file)

			if cmd.Flag(defaultValFlag).Changed {
				item, err = repo.Find(key, defaultVal)
			} else {
				item, err = repo.Get(key)
			}

			if err != nil {
				return err
			}

			cmd.Print(item.Val)

			return nil
		},
	}

	defaultVal = cmd.Flags().StringP(defaultValFlag, "d", "", "Default value")

	return cmd
}

func init() {
	rootCmd.AddCommand(newGetCmd())
}
