package github

import (
	"github.com/everettraven/synkr/pkg/plugins"
	"go.starlark.net/starlark"
)

func init() {
	plugins.Register(New())
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
