package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/everettraven/synkr/pkg/builtins"
	"github.com/everettraven/synkr/pkg/engine"
	"github.com/everettraven/synkr/pkg/printers"
	"github.com/spf13/cobra"
	"go.starlark.net/lib/time"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

func NewSynkrCommand() *cobra.Command {
	eng := &engine.Engine{}
	var configFile string
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "synkr",
		Short: "synkr is an engine for syncing work items based on a Starlark configuration",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), eng, configFile, outputFormat)
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", defaultConfigPath(), "configures the Starlark file to be processed for configuration. Defaults to $HOME/.config/synkr.star if possible to get your home directory. Otherwise it uses synkr.star in the current directory.")
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "markdown", "configures the output format. Allowed values are [markdown, json]")

	return cmd
}

func run(ctx context.Context, eng *engine.Engine, configFile, output string) error {
	thread, err := configureEngine(eng, configFile, output)
	if err != nil {
		return fmt.Errorf("configuring engine: %w", err)
	}

	return eng.Run(ctx, thread)
}

func configureEngine(eng *engine.Engine, configFile, output string) (*starlark.Thread, error) {
	switch output {
	case "markdown":
		eng.SetPrinter(&printers.Markdown{})
	case "json":
		eng.SetPrinter(&printers.JSON{})
	default:
		return nil, fmt.Errorf("unknown output format %q", output)
	}

	globals := starlark.StringDict{}
	starlark.Universe["time"] = time.Module

	builtins.Github(globals, eng)

	thread := &starlark.Thread{Name: "main"}

	_, err := starlark.ExecFileOptions(&syntax.FileOptions{TopLevelControl: true}, thread, configFile, nil, globals)

	return thread, err
}

// defaultConfigPath returns the default configuration file path that should be used.
// If we can get the users home directory we default to $HOME/.config/synkr.star, otherwise
// we default to synkr.star
func defaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "synkr.star"
	}

	return filepath.Join(homeDir, ".config", "synkr.star")
}
