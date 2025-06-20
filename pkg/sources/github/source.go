package github

import (
	"context"
	"fmt"

	"go.starlark.net/starlark"
)

type GitHub struct {
	filters []starlark.Callable
	org     string
	repo    string
}

func New(org, repo string, filters ...starlark.Callable) *GitHub {
	return &GitHub{
		org:     org,
		repo:    repo,
		filters: filters,
	}
}

func (g *GitHub) Name() string {
	return "GitHub"
}

func (g *GitHub) Project() string {
	return fmt.Sprintf("%s/%s", g.org, g.repo)
}

func (g *GitHub) Fetch(ctx context.Context, thread *starlark.Thread) ([]any, error) {
	items := []RepoItem{}
	issues, err := getIssuesForRepo(ctx, g.org, g.repo)
	if err != nil {
		return nil, fmt.Errorf("fetching issues for GitHub repository %s/%s : %w", g.org, g.repo, err)
	}
	items = append(items, issues...)

	pullRequests, err := getPullRequestsForRepo(ctx, g.org, g.repo)
	if err != nil {
		return nil, fmt.Errorf("fetching pull requests for GitHub repository %s/%s : %w", g.org, g.repo, err)
	}
	items = append(items, pullRequests...)

	return g.filterItems(thread, items...)
}

func (g *GitHub) filterItems(thread *starlark.Thread, items ...RepoItem) ([]any, error) {
	outItems := []any{}

	for _, item := range items {
		include := true
		for _, filter := range g.filters {
			out, err := starlark.Call(thread, filter, starlark.Tuple{repoItemToStarlarkDict(item)}, nil)
			if err != nil {
				return nil, err
			}

			if !out.Truth() {
				include = false
				break
			}
		}

		if !include {
			continue
		}

		outItems = append(outItems, item)
	}

	return outItems, nil
}
