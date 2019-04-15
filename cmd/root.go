package cmd

import (
	"github.com/filkra/ilias/cmd/submissions"
	"github.com/spf13/cobra"
	"os"
)

var rootCommand = &cobra.Command{
	Use:   "ilias",
	Short: "A simple command line interface for managing ILIAS",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	rootCommand.AddCommand(submissions.RootCommand)
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
