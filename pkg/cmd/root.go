package cmd

import (
	"github.com/everettraven/wranglr/pkg/runner"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	runOpts := &runner.Options{}

	cmd := &cobra.Command{
		Use:   "wranglr",
		Short: "wranglr is an engine for wrangling together work items based on a Starlark configuration",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runOpts.Run(cmd.Context())
		},
	}

	cmd.Flags().StringVarP(&runOpts.ConfigFile, "config", "c", runner.DefaultConfigPath(), "configures the Starlark file to be processed for configuration. Defaults to $HOME/.config/wranglr.star if possible to get your home directory. Otherwise it uses wranglr.star in the current directory.")
	cmd.Flags().StringVarP(&runOpts.OutputFormat, "output", "o", "interactive", "configures the output format. Allowed values are [json, interactive]")

	return cmd
}
