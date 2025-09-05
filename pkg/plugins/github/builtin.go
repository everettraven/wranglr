package github

import (
	"github.com/everettraven/synkr/pkg/builtins"
	"github.com/google/go-github/v71/github"
	"go.starlark.net/starlark"
)

func GithubBuiltinFunc(sourcer *Sourcer) builtins.BuiltinFunc {
	return func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var filters *starlark.List
		var priorities *starlark.List
		var status starlark.Callable
		var org starlark.String
		var repo starlark.String

		// GitHub issue list filters
		var milestone starlark.String
		var state starlark.String
		var assignee starlark.String
		var creator starlark.String
		var mentioned starlark.String
		var labels *starlark.List
		var sort starlark.String
		var direction starlark.String

		err := starlark.UnpackArgs("github", args, kwargs,
			"org", &org,
			"repo", &repo,
			"milestone?", &milestone,
			"state?", &state,
			"assignee?", &assignee,
			"creator?", &creator,
			"mentioned?", &mentioned,
			"labels?", &labels,
			"sort?", &sort,
			"direction?", &direction,
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

		listOpts := &github.IssueListByRepoOptions{
			Milestone: milestone.GoString(),
			State:     state.GoString(),
			Assignee:  assignee.GoString(),
			Creator:   creator.GoString(),
			Mentioned: mentioned.GoString(),
			Labels:    builtins.TypeFromStarlarkList[string](labels),
			Sort:      sort.GoString(),
			Direction: direction.GoString(),
		}

		ghSource := NewSource(org.GoString(), repo.GoString(), filterCallables, priorityCallables, status, listOpts)
		sourcer.AddSource(ghSource)

		return starlark.None, nil
	}
}

func ptr[T any](in T) *T {
	return &in
}
