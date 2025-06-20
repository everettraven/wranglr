package cmd

import (
	"context"
	"fmt"

	"github.com/everettraven/synkr/pkg/builtins"
	"github.com/everettraven/synkr/pkg/engine"
	"github.com/spf13/cobra"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

func NewSynkrCommand() *cobra.Command {
	eng := &engine.Engine{}
	var configFile string

	cmd := &cobra.Command{
		Use:   "synkr [-c configFile]",
		Short: "synkr is an engine for syncing work items based on a Starlark configuration",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), eng, configFile)
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", "synkr.star", "configures the Starlark file to be processed for configuration")

	return cmd
}

func run(ctx context.Context, eng *engine.Engine, configFile string) error {
	thread, err := configureEngine(eng, configFile)
	if err != nil {
		return fmt.Errorf("configuring engine: %w", err)
	}

	return eng.Run(ctx, thread)
}

func configureEngine(eng *engine.Engine, configFile string) (*starlark.Thread, error) {
	globals := starlark.StringDict{}

	builtins.Github(globals, eng)

	thread := &starlark.Thread{Name: "main"}

	_, err := starlark.ExecFileOptions(&syntax.FileOptions{}, thread, configFile, nil, globals)

	return thread, err
}
