package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cli/cli/v2/pkg/search"
	"github.com/everettraven/synkr/pkg/plugins"
	"go.starlark.net/starlark"
)

type Source struct {
	filters    []starlark.Callable
	priorities []starlark.Callable
	status     starlark.Callable
	searcher   search.Searcher
	query      search.Query
}

func NewSource(host string, filters, priorities []starlark.Callable, status starlark.Callable, query search.Query) *Source {
	return &Source{
		filters:    filters,
		priorities: priorities,
		status:     status,
		searcher:   newGithubSearcher(host),
		query:      query,
	}
}

func (g *Source) Name() string {
	return fmt.Sprintf("github/%s", g.Project())
}

func (g *Source) Project() string {
	return g.query.Qualifiers.Repo[0]
}

func (g *Source) Fetch(ctx context.Context, thread *starlark.Thread) (*plugins.SourceResult, error) {
	items := []RepoItem{}
	issues, err := getIssuesForRepo(ctx, g.searcher, g.query)
	if err != nil {
		return nil, fmt.Errorf("fetching issues for GitHub repository %q : %w", g.query.Qualifiers.Repo[0], err)
	}
	items = append(items, issues...)

	items, err = g.filterItems(thread, items...)
	if err != nil {
		return nil, fmt.Errorf("filtering items for GitHub repository %q: %w", g.query.Qualifiers.Repo[0], err)
	}

	items, err = g.setPriority(thread, items...)
	if err != nil {
		return nil, fmt.Errorf("setting item priorities for GitHub repository %q: %w", g.query.Qualifiers.Repo[0], err)
	}

	if g.status != nil {
		items, err = g.setStatus(thread, items...)
		if err != nil {
			return nil, fmt.Errorf("setting item statuses for GitHub repository %q: %w", g.query.Qualifiers.Repo[0], err)
		}
	}

	return &plugins.SourceResult{
		Source:  "GitHub",
		Project: g.Project(),
		Items:   repoItemSliceToSourceEntrySlice(items...),
	}, nil
}

func (g *Source) filterItems(thread *starlark.Thread, items ...RepoItem) ([]RepoItem, error) {
	if len(g.filters) == 0 {
		return items, nil
	}

	outItems := []RepoItem{}

	for _, item := range items {
		for _, filter := range g.filters {
			out, err := starlark.Call(thread, filter, starlark.Tuple{repoItemToStarlarkDict(item)}, nil)
			if err != nil {
				return nil, err
			}

			if out.Truth() {
				outItems = append(outItems, item)
			}
		}
	}

	return outItems, nil
}

func (g *Source) setPriority(thread *starlark.Thread, items ...RepoItem) ([]RepoItem, error) {
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

func (g *Source) setStatus(thread *starlark.Thread, items ...RepoItem) ([]RepoItem, error) {
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

func newGithubSearcher(host string) search.Searcher {
	return search.NewSearcher(http.DefaultClient, host)
}

func repoItemSliceToSourceEntrySlice(items ...RepoItem) []plugins.SourceEntry {
	out := []plugins.SourceEntry{}

	for _, item := range items {
		out = append(out, item)
	}

	return out
}

type RepoItemType string

const (
	RepoItemTypeIssue       RepoItemType = "Issue"
	RepoItemTypePullRequest RepoItemType = "PullRequest"
)

// TODO: Expand here to capture more things that
// may be important to filter on
type RepoItem struct {
	ID        string       `json:"id"`
	URL       string       `json:"url"`
	Author    string       `json:"author"`
	Labels    []string     `json:"labels"`
	Type      RepoItemType `json:"type"`
	Assignees []string     `json:"assignees"`
	Title     string       `json:"title"`
	Body      string       `json:"body"`
	State     string       `json:"state"`
	Priority  int          `json:"priority"`
	Status    string       `json:"status"`
	Created   string       `json:"created"`
	Updated   string       `json:"updated"`
	Comments  int          `json:"comments"`

	// Only populated on PullRequests
	RequestedReviewers []string `json:"requestedReviewers"`
	Draft              bool     `json:"draft"`
}

func (ri RepoItem) Identifier() string {
	return ri.ID
}

func repoItemToStarlarkDict(item RepoItem) *starlark.Dict {
	dict := &starlark.Dict{}

	// TODO: handle errors when setting keys
	_ = dict.SetKey(starlark.String("author"), starlark.String(item.Author))
	_ = dict.SetKey(starlark.String("type"), starlark.String(item.Type))
	_ = dict.SetKey(starlark.String("title"), starlark.String(item.Title))
	_ = dict.SetKey(starlark.String("body"), starlark.String(item.Body))
	_ = dict.SetKey(starlark.String("state"), starlark.String(item.State))
	_ = dict.SetKey(starlark.String("labels"), starlark.NewList(stringArrayToStarlarkValueArray(item.Labels...)))
	_ = dict.SetKey(starlark.String("assignees"), starlark.NewList(stringArrayToStarlarkValueArray(item.Assignees...)))
	_ = dict.SetKey(starlark.String("created"), starlark.String(item.Created))
	_ = dict.SetKey(starlark.String("updated"), starlark.String(item.Updated))
	_ = dict.SetKey(starlark.String("comments"), starlark.MakeInt(item.Comments))

	if item.Type == RepoItemTypePullRequest {
		_ = dict.SetKey(starlark.String("draft"), starlark.Bool(item.Draft))
		_ = dict.SetKey(starlark.String("requestedReviewers"), starlark.NewList(stringArrayToStarlarkValueArray(item.RequestedReviewers...)))
	}

	return dict
}

func stringArrayToStarlarkValueArray(in ...string) []starlark.Value {
	out := []starlark.Value{}
	for _, i := range in {
		out = append(out, starlark.String(i))
	}

	return out
}

func getIssuesForRepo(ctx context.Context, searcher search.Searcher, query search.Query) ([]RepoItem, error) {
	issues, err := searcher.Issues(query)
	if err != nil {
		return nil, fmt.Errorf("fetching issues: %v", err)
	}

	out := []RepoItem{}
	for _, issue := range issues.Items {
		item := RepoItem{
			Type:     RepoItemTypeIssue,
			ID:       issue.ID,
			URL:      issue.URL,
			Author:   issue.Author.Login,
			Title:    issue.Title,
			Body:     issue.Body,
			State:    issue.State(),
			Created:  issue.CreatedAt.String(),
			Updated:  issue.UpdatedAt.String(),
			Comments: issue.CommentsCount,
		}

		if issue.IsPullRequest() {
			item.Type = RepoItemTypePullRequest
			item.URL = issue.PullRequest.URL
		}

		if issue.Assignees != nil {
			assignees := []string{}
			for _, assignee := range issue.Assignees {
				assignees = append(assignees, assignee.Login)
			}

			item.Assignees = assignees
		}

		if issue.Labels != nil {
			labels := []string{}
			for _, label := range issue.Labels {
				labels = append(labels, label.Name)
			}

			item.Labels = labels
		}

		out = append(out, item)
	}

	return out, nil
}
