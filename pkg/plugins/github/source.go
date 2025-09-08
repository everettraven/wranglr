package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

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

func (g *Source) Fetch(ctx context.Context, thread *starlark.Thread, resultChan chan plugins.SourceEntry) error {
	queueChan := make(chan *RepoItem)

	fetchWaitGroup := sync.WaitGroup{}
	fetchWaitGroup.Add(1)
	errs := []error{}
	go func() {
		err := g.getIssuesForRepo(queueChan)
		if err != nil {
			errs = append(errs, fmt.Errorf("fetching issues for GitHub repository %q : %w", g.query.Qualifiers.Repo[0], err))
		}
		fetchWaitGroup.Done()
	}()

	processGroup := sync.WaitGroup{}
	processGroup.Add(1)
	go func() {
		for {
			shouldExit := false
			select {
			case <-ctx.Done():
				shouldExit = true
			case item, ok := <-queueChan:
				if !ok {
					shouldExit = true
					break
				}

				// process items
				include, err := g.checkFilters(thread, item)
				if err != nil {
					errs = append(errs, err)
				}

				if !include {
					break
				}

				err = g.setStatus(thread, item)
				if err != nil {
					errs = append(errs, err)
					break
				}

				err = g.setPriority(thread, item)
				if err != nil {
					errs = append(errs, err)
					break
				}

				resultChan <- item
			}

			if shouldExit {
				break
			}
		}
		processGroup.Done()
	}()

	fetchWaitGroup.Wait()
	close(queueChan)

	processGroup.Wait()

	return errors.Join(errs...)
}

func (g *Source) checkFilters(thread *starlark.Thread, item *RepoItem) (bool, error) {
	for _, filterFunc := range g.filters {
		val, err := starlark.Call(thread, filterFunc, starlark.Tuple{repoItemToStarlarkDict(item)}, nil)
		if err != nil {
			return false, fmt.Errorf("calling filter function %q: %w", filterFunc.Name(), err)
		}

		if !val.Truth() {
			return false, nil
		}
	}

	return true, nil
}

func (g *Source) setPriority(thread *starlark.Thread, item *RepoItem) error {
	itemScore := 0
	for _, priorityFunc := range g.priorities {
		val, err := starlark.Call(thread, priorityFunc, starlark.Tuple{repoItemToStarlarkDict(item)}, nil)
		if err != nil {
			return fmt.Errorf("calling priority function %q: %w", priorityFunc.Name(), err)
		}

		score := 0
		err = starlark.AsInt(val, &score)
		if err != nil {
			return fmt.Errorf("could not use return value of priority function %q as an integer: %w", priorityFunc.Name(), err)
		}

		itemScore += score
	}

	item.Priority = itemScore
	return nil
}

func (g *Source) setStatus(thread *starlark.Thread, item *RepoItem) error {
	if g.status == nil {
		return nil
	}

	val, err := starlark.Call(thread, g.status, starlark.Tuple{repoItemToStarlarkDict(item)}, nil)
	if err != nil {
		return fmt.Errorf("calling status function %q: %w", g.status.Name(), err)
	}

	status, _ := starlark.AsString(val)
	item.Status = status

	return nil
}

func newGithubSearcher(host string) search.Searcher {
	return search.NewSearcher(http.DefaultClient, host)
}

type RepoItemType string

const (
	RepoItemTypeIssue       RepoItemType = "Issue"
	RepoItemTypePullRequest RepoItemType = "PullRequest"
)

// TODO: Expand here to capture more things that
// may be important to filter on
type RepoItem struct {
	ID         string       `json:"id"`
	URL        string       `json:"url"`
	Author     string       `json:"author"`
	Labels     []string     `json:"labels"`
	Type       RepoItemType `json:"type"`
	Assignees  []string     `json:"assignees"`
	Title      string       `json:"title"`
	Body       string       `json:"body"`
	State      string       `json:"state"`
	Priority   int          `json:"priority"`
	Status     string       `json:"status"`
	Created    string       `json:"created"`
	Updated    string       `json:"updated"`
	Comments   int          `json:"comments"`
	Project    string       `json:"project"`
	SourceName string       `json:"source"`

	// Only populated on PullRequests
	RequestedReviewers []string `json:"requestedReviewers"`
	Draft              bool     `json:"draft"`
}

func (ri RepoItem) Identifier() string {
	return ri.ID
}

func (ri RepoItem) Source() string {
	return fmt.Sprintf("%s/%s", ri.SourceName, ri.Project)
}

func repoItemToStarlarkDict(item *RepoItem) *starlark.Dict {
	dict := &starlark.Dict{}

	if item == nil {
		return dict
	}

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

func (g *Source) getIssuesForRepo(queue chan *RepoItem) error {
	issues, err := g.searcher.Issues(g.query)
	if err != nil {
		return fmt.Errorf("fetching issues: %v", err)
	}

	for _, issue := range issues.Items {
		item := &RepoItem{
			Type:       RepoItemTypeIssue,
			ID:         issue.ID,
			URL:        issue.URL,
			Author:     issue.Author.Login,
			Title:      issue.Title,
			Body:       issue.Body,
			State:      issue.State(),
			Created:    issue.CreatedAt.String(),
			Updated:    issue.UpdatedAt.String(),
			Comments:   issue.CommentsCount,
			SourceName: "github",
			Project:    g.Project(),
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

		queue <- item
	}

	return nil
}
