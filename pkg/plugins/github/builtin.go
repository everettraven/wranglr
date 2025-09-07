package github

import (
	"github.com/cli/cli/v2/pkg/search"
	"github.com/everettraven/synkr/pkg/builtins"
	"go.starlark.net/starlark"
)

func GithubBuiltinFunc(sourcer *Sourcer) builtins.BuiltinFunc {
	return func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var priorities *starlark.List
		var status starlark.Callable

		var host starlark.String
		var repo starlark.String

		// GitHub search qualifiers
		var assignee starlark.String
		var author starlark.String
		var closed starlark.String
		var commenter starlark.String
		var comments starlark.String
		var created starlark.String
		var draft *starlark.Bool
		var extension starlark.String
		var filename starlark.String
		var in *starlark.List
		var involves starlark.String
		var is *starlark.List
		var labels *starlark.List
		var language starlark.String
		var mentions starlark.String
		var merged starlark.String
		var milestone starlark.String
		var no *starlark.List
		var path starlark.String
		var review starlark.String
		var reviewRequested starlark.String
		var reviewedBy starlark.String
		var state starlark.String
		var team starlark.String
		var teamReviewRequested starlark.String
		var updated starlark.String
		var user *starlark.List

		// GitHub search top-level configurations
		var sort starlark.String
		var order starlark.String
		var limit starlark.Int

		err := starlark.UnpackArgs("github", args, kwargs,
			"host?", &host,
			"repo", &repo,
			"assignee?", &assignee,
			"author?", &author,
			"closed?", &closed,
			"commenter", &commenter,
			"comments", &comments,
			"created?", &created,
			"draft?", &draft,
			"extension?", &extension,
			"filename?", &filename,
			"in?", &in,
			"involves?", &involves,
			"is?", &is,
			"labels?", &labels,
			"language?", &language,
			"mentions?", &mentions,
			"merged?", &merged,
			"milestone?", &milestone,
			"no?", &no,
			"path?", &path,
			"review?", &review,
			"review_requested?", &reviewRequested,
			"reviewed_by?", &reviewedBy,
			"state?", &state,
			"team?", &team,
			"team_review_requested?", &teamReviewRequested,
			"updated?", &updated,
			"user?", &user,
			"sort?", &sort,
			"order?", &order,
			"limit?", &limit,
			"priorities?", &priorities,
			"status?", &status,
		)
		if err != nil {
			return nil, err
		}

		priorityCallables := []starlark.Callable{}
		if priorities != nil {
			priorityCallables = builtins.TypeFromStarlarkList[starlark.Callable](priorities)
		}

		limitValue := 100
		if limit.BigInt().Int64() > 0 {
			limitValue = int(limit.BigInt().Int64())
		}

		var draftValue *bool
		if draft != nil {
			draftValue = ptr(bool(draft.Truth()))
		}

		query := search.Query{
			Limit: limitValue,
			Kind:  search.KindIssues,
			Order: order.GoString(),
			Sort:  sort.GoString(),
			Qualifiers: search.Qualifiers{
				Assignee:            assignee.GoString(),
				Author:              author.GoString(),
				Closed:              closed.GoString(),
				Commenter:           commenter.GoString(),
				Comments:            comments.GoString(),
				Created:             created.GoString(),
				Draft:               draftValue,
				Extension:           extension.GoString(),
				Filename:            filename.GoString(),
				In:                  builtins.TypeFromStarlarkList[string](in),
				Involves:            involves.GoString(),
				Is:                  builtins.TypeFromStarlarkList[string](is),
				Label:               builtins.TypeFromStarlarkList[string](labels),
				Language:            language.GoString(),
				Merged:              merged.GoString(),
				Mentions:            mentions.GoString(),
				Milestone:           milestone.GoString(),
				Repo:                []string{repo.GoString()},
				State:               state.GoString(),
				No:                  builtins.TypeFromStarlarkList[string](no),
				Path:                path.GoString(),
				ReviewRequested:     reviewRequested.GoString(),
				ReviewedBy:          reviewedBy.GoString(),
				Team:                team.GoString(),
				TeamReviewRequested: teamReviewRequested.GoString(),
				Updated:             updated.GoString(),
				User:                builtins.TypeFromStarlarkList[string](user),
			},
		}

		hostValue := "github.com"
		if host.GoString() != "" {
			hostValue = host.GoString()
		}

		ghSource := NewSource(hostValue, priorityCallables, status, query)
		sourcer.AddSource(ghSource)

		return starlark.None, nil
	}
}

func ptr[T any](in T) *T {
	return &in
}
