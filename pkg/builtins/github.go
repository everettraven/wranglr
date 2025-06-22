package builtins

import (
	"github.com/everettraven/synkr/pkg/engine"
	"github.com/everettraven/synkr/pkg/sources/github"
	"go.starlark.net/starlark"
)

func Github(global starlark.StringDict, eng *engine.Engine) {
	global["github"] = starlark.NewBuiltin("github", githubBuiltinFunc(eng))
}

func githubBuiltinFunc(eng *engine.Engine) BuiltinFunc {
	return func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var filters *starlark.List
		var priorities *starlark.List
		var org starlark.String
		var repo starlark.String

		err := starlark.UnpackArgs("github", args, kwargs, "org", &org, "repo", &repo, "filters?", &filters, "priorities?", &priorities)
		if err != nil {
			return nil, err
		}

		filterCallables := []starlark.Callable{}
		if filters != nil {
			filterCallables = callablesFromList(filters)
		}

		priorityCallables := []starlark.Callable{}
		if priorities != nil {
			priorityCallables = callablesFromList(priorities)
		}

		ghSource := github.New(org.GoString(), repo.GoString(), filterCallables, priorityCallables)
		eng.AddSource(ghSource)

		return starlark.None, nil
	}
}

func callablesFromList(list *starlark.List) []starlark.Callable {
	callables := []starlark.Callable{}

	if list != nil {
		for v := range list.Elements() {
			if callable, ok := v.(starlark.Callable); ok {
				callables = append(callables, callable)
			}
		}
	}

	return callables
}
