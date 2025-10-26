package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"go.starlark.net/lib/time"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"

	"github.com/everettraven/wranglr/pkg/modules"
	"github.com/everettraven/wranglr/pkg/modules/github"
	"github.com/everettraven/wranglr/pkg/modules/jira"
	"github.com/everettraven/wranglr/pkg/modules/wranglr"
)

type Options struct {
	ConfigFile   string
	OutputFormat string
}

func (o *Options) Run(ctx context.Context) error {
	// register all modules
	err := modules.Register(github.New())
	if err != nil {
		return err
	}

	err = modules.Register(jira.New())
	if err != nil {
		return err
	}

	err = modules.Register(wranglr.New(o.OutputFormat))
	if err != nil {
		return err
	}

	// Do actual things
	_, err = configureThread(o.ConfigFile, o.OutputFormat)
	if err != nil {
		return fmt.Errorf("configuring thread: %w", err)
	}
	return nil
}

func configureThread(configFile string, output string) (*starlark.Thread, error) {
	globals := starlark.StringDict{}
	starlark.Universe["time"] = time.Module

	for name, module := range modules.Modules() {
		globals[name] = module
	}

	thread := &starlark.Thread{Name: "main"}

	_, err := starlark.ExecFileOptions(
		&syntax.FileOptions{
			TopLevelControl: true,
			GlobalReassign:  true,
		},
		thread,
		configFile,
		nil,
		globals,
	)

	return thread, err
}

// DefaultConfigPath returns the default configuration file path that should be used.
// If we can get the users home directory we default to $HOME/.config/wranglr.star, otherwise
// we default to wranglr.star
func DefaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "wranglr.star"
	}

	return filepath.Join(homeDir, ".config", "wranglr.star")
}
