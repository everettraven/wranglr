package github

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/v71/github"
	"go.starlark.net/starlark"
)

type GitHub struct {
	filters         []starlark.Callable
	priorities      []starlark.Callable
	status          starlark.Callable
	org             string
	repo            string
	client          *github.Client
	includeMentions bool
}

func New(org, repo string, filters []starlark.Callable, priorities []starlark.Callable, status starlark.Callable, includeMentions bool) *GitHub {
	return &GitHub{
		org:             org,
		repo:            repo,
		filters:         filters,
		client:          newGithubClient(),
		priorities:      priorities,
		status:          status,
		includeMentions: includeMentions,
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
	issues, err := getIssuesForRepo(ctx, g.client, g.org, g.repo, g.includeMentions)
	if err != nil {
		return nil, fmt.Errorf("fetching issues for GitHub repository %s/%s : %w", g.org, g.repo, err)
	}
	items = append(items, issues...)

	pullRequests, err := getPullRequestsForRepo(ctx, g.client, g.org, g.repo, g.includeMentions)
	if err != nil {
		return nil, fmt.Errorf("fetching pull requests for GitHub repository %s/%s : %w", g.org, g.repo, err)
	}
	items = append(items, pullRequests...)

	items, err = g.filterItems(thread, items...)
	if err != nil {
		return nil, fmt.Errorf("filtering items for GitHub repository %s/%s: %w", g.org, g.repo, err)
	}

	items, err = g.setPriority(thread, items...)
	if err != nil {
		return nil, fmt.Errorf("setting item priorities for GitHub repository %s/%s: %w", g.org, g.repo, err)
	}

	if g.status != nil {
		items, err = g.setStatus(thread, items...)
		if err != nil {
			return nil, fmt.Errorf("setting item statuses for GitHub repository %s/%s: %w", g.org, g.repo, err)
		}
	}

	return repoItemSliceToAnySlice(items...), nil
}

func (g *GitHub) filterItems(thread *starlark.Thread, items ...RepoItem) ([]RepoItem, error) {
	outItems := []RepoItem{}

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

func (g *GitHub) setPriority(thread *starlark.Thread, items ...RepoItem) ([]RepoItem, error) {
	out := []RepoItem{}

	for _, item := range items {
		itemScore := 0
		for _, priorityFunc := range g.priorities {
			val, err := starlark.Call(thread, priorityFunc, starlark.Tuple{repoItemToStarlarkDict(item)}, nil)
			if err != nil {
				return nil, fmt.Errorf("calling priority function %q: %w", priorityFunc.Name(), err)
			}

			score := 0
			err = starlark.AsInt(val, &score)
			if err != nil {
				return nil, fmt.Errorf("could not use return value of priority function %q as an integer: %w", priorityFunc.Name(), err)
			}

			itemScore += score
		}

		item.Priority = itemScore
		out = append(out, item)
	}

	return out, nil
}

func (g *GitHub) setStatus(thread *starlark.Thread, items ...RepoItem) ([]RepoItem, error) {
	out := []RepoItem{}

	for _, item := range items {
		val, err := starlark.Call(thread, g.status, starlark.Tuple{repoItemToStarlarkDict(item)}, nil)
		if err != nil {
			return nil, fmt.Errorf("calling status function %q: %w", g.status.Name(), err)
		}

		status, _ := starlark.AsString(val)
		item.Status = status
		out = append(out, item)
	}

	return out, nil
}

func newGithubClient() *github.Client {
	token := os.Getenv("SYNKR_GITHUB_TOKEN")
	if token == "" {
		return github.NewClient(http.DefaultClient)
	}

	return github.NewClient(http.DefaultClient).WithAuthToken(token)
}

func repoItemSliceToAnySlice(items ...RepoItem) []any {
	out := []any{}

	for _, item := range items {
		out = append(out, item)
	}

	return out
}
