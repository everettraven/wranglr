package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/everettraven/synkr/pkg/engine"
	"github.com/everettraven/synkr/pkg/plugins"
	"github.com/everettraven/synkr/pkg/printers"
	"github.com/spf13/cobra"
	"go.starlark.net/lib/time"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"

	_ "github.com/everettraven/synkr/pkg/plugins/registration"
)

func NewSynkrCommand() *cobra.Command {
	var configFile string
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "synkr",
		Short: "synkr is an engine for syncing work items based on a Starlark configuration",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), configFile, outputFormat)
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", defaultConfigPath(), "configures the Starlark file to be processed for configuration. Defaults to $HOME/.config/synkr.star if possible to get your home directory. Otherwise it uses synkr.star in the current directory.")
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "markdown", "configures the output format. Allowed values are [markdown, json, web]")

	return cmd
}

func run(ctx context.Context, configFile, output string) error {
	plugins := plugins.Plugins()
	thread, err := configureThread(configFile, plugins...)
	if err != nil {
		return fmt.Errorf("configuring thread: %w", err)
	}

	eng := engine.New(plugins...)

	results, err := eng.Run(ctx, thread)
	if err != nil {
		return fmt.Errorf("running engine: %w", err)
	}

	return printResults(output, results...)
}

func printResults(output string, results ...plugins.SourceResult) error {
	switch output {
	case "json":
		out := &printers.JSON{}
		return out.Print(results...)
	case "markdown":
		out := &printers.Markdown{}
		return out.Print(results...)
	case "web":
		out := &printers.Web{}
		return out.Print(results...)
	default:
		return fmt.Errorf("unknown output format %q", output)
	}
}

func configureThread(configFile string, plugins ...plugins.Plugin) (*starlark.Thread, error) {
	globals := starlark.StringDict{}
	starlark.Universe["time"] = time.Module

	for _, plugin := range plugins {
		plugin.RegisterBuiltins(globals)
	}

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
