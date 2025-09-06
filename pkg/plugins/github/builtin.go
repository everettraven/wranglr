package github

import (
	"github.com/cli/cli/v2/pkg/search"
	"github.com/everettraven/synkr/pkg/builtins"
	"go.starlark.net/starlark"
)

func GithubBuiltinFunc(sourcer *Sourcer) builtins.BuiltinFunc {
	return func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var filters *starlark.List
		var priorities *starlark.List
		var status starlark.Callable

		var host starlark.String
		var repo starlark.String

		// GitHub issue list filters
		// TODO: expand this and/or make this the return value of a callable
		var state starlark.String
		var assignee starlark.String
		var creator starlark.String
		var mentioned starlark.String
		var labels *starlark.List
		var sort starlark.String
		var direction starlark.String
		var limit starlark.Int

		err := starlark.UnpackArgs("github", args, kwargs,
			"host?", &host,
			"repo", &repo,
			"state?", &state,
			"assignee?", &assignee,
			"creator?", &creator,
			"mentioned?", &mentioned,
			"labels?", &labels,
			"sort?", &sort,
			"direction?", &direction,
			"limit?", &limit,
			"filters?", &filters,
			"priorities?", &priorities,
			"status?", &status,
		)
		if err != nil {
			return nil, err
		}

		filterCallables := []starlark.Callable{}
		if filters != nil {
			filterCallables = builtins.TypeFromStarlarkList[starlark.Callable](filters)
		}

		priorityCallables := []starlark.Callable{}
		if priorities != nil {
			priorityCallables = builtins.TypeFromStarlarkList[starlark.Callable](priorities)
		}

		limitValue := 100
		if limit.BigInt().Int64() > 0 {
			limitValue = int(limit.BigInt().Int64())
		}

		query := search.Query{
			Limit: limitValue,
			Kind:  search.KindIssues,
			Order: direction.GoString(),
			Sort:  sort.GoString(),
			Qualifiers: search.Qualifiers{
				State:    state.GoString(),
				Assignee: assignee.GoString(),
				Author:   creator.GoString(),
				Mentions: mentioned.GoString(),
				Label:    builtins.TypeFromStarlarkList[string](labels),
				Repo:     []string{repo.GoString()},
			},
		}

		hostValue := "github.com"
		if host.GoString() != "" {
			hostValue = host.GoString()
		}

		ghSource := NewSource(hostValue, filterCallables, priorityCallables, status, query)
		sourcer.AddSource(ghSource)

		return starlark.None, nil
	}
}
