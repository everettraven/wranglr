package github

import (
	"fmt"

	"github.com/everettraven/synkr/pkg/plugins"
	"go.starlark.net/starlark"
)

func init() {
	err := plugins.Register(New())
	if err != nil {
		panic(fmt.Errorf("registering github plugin: %w", err))
	}
}

func New() plugins.Plugin {
	sourcer := &Sourcer{
		sources: make([]plugins.Source, 0),
	}
	return plugins.Plugin{
		Sourcer: sourcer,
		Builtins: map[string]*starlark.Builtin{
			"github": starlark.NewBuiltin("github", GithubBuiltinFunc(sourcer)),
		},
		Name: "github",
	}
}
