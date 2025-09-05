package plugins

import (
	"context"
	"fmt"
	"slices"

	"go.starlark.net/starlark"
)

type SourceResult struct {
	Source  string `json:"source"`
	Project string `json:"project"`
	Items   []SourceEntry  `json:"items"`
}

type SourceEntry interface {
	Identifier() string
}

type Source interface {
	Fetch(context.Context, *starlark.Thread) (*SourceResult, error)
	Name() string
}

type Sourcer interface {
	Sources() []Source
}

type Plugin struct {
	Sourcer
	Builtins map[string]*starlark.Builtin
	Name string
}

func (p Plugin) RegisterBuiltins(global starlark.StringDict) {
	for key, builtin := range p.Builtins {
		global[key] = builtin
	}
}

var plugins []Plugin

func Register(plugin Plugin) error {
	if slices.ContainsFunc(plugins, func(e Plugin) bool {
		return e.Name == plugin.Name
	}){
		return fmt.Errorf("plugin %q is already registered", plugin.Name)
	}

	plugins = append(plugins, plugin)
	return nil
}

func Plugins() []Plugin {
	return plugins
}
