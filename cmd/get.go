package cmd

import (
	"fmt"
	"github.com/oxio/kv/internal/kv"
	"github.com/oxio/kv/internal/parser"
	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	var defaultVal *string
	var skipMissingFiles *bool

	cmd := &cobra.Command{
		Use:   "get <file1> [<file2> <file3> ...] <key> [--default|-d value] [--skip-missing-files|-m]",
		Short: "Gets a value from the key-value file",
		RunE: func(cmd *cobra.Command, args []string) error {
			var item *parser.Item
			var err error

			if len(args) < 2 {
				return fmt.Errorf("not enough arguments")
			}

			files := args[:len(args)-1]
			key := args[len(args)-1]

			for _, file := range files {
				repo := kv.NewKvRepo(file, *skipMissingFiles)

				if cmd.Flag(defaultValFlag).Changed {
					item, err = repo.Find(key, defaultVal)
				} else {
					item, err = repo.Get(key)
				}
			}

			if err != nil {
				return err
			}

			cmd.Print(item.Val)

			return nil
		},
	}

	defaultVal = cmd.Flags().StringP(
		defaultValFlag,
		defaultValShortFlag,
		"",
		"This value will be returned if the key is not found in the provided file(s).",
	)
	skipMissingFiles = cmd.Flags().BoolP(
		skipMissingFilesFlag,
		skipMissingFilesShortFlag,
		false,
		"Do not issue \"no such file or directory\" error on missing or inaccessible files. Should only be used"+
			" with multiple files or in combination with the \"--default\" flag.",
	)

	return cmd
}

func init() {
	rootCmd.AddCommand(newGetCmd())
}
