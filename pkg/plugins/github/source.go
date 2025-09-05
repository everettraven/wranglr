package github

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/everettraven/synkr/pkg/plugins"
	"github.com/google/go-github/v71/github"
	"go.starlark.net/starlark"
)

type Source struct {
	filters          []starlark.Callable
	priorities       []starlark.Callable
	status           starlark.Callable
	org              string
	repo             string
	client           *github.Client
	issueListOptions *github.IssueListByRepoOptions
}

func NewSource(org, repo string, filters, priorities []starlark.Callable, status starlark.Callable, issueListOptions *github.IssueListByRepoOptions) *Source {
	return &Source{
		org:              org,
		repo:             repo,
		filters:          filters,
		// TODO: should this just hit the 'gh' cli instead of hitting GH API?
		client:           newGithubClient(),
		priorities:       priorities,
		status:           status,
		issueListOptions: issueListOptions,
	}
}

func (g *Source) Name() string {
	return fmt.Sprintf("github/%s", g.Project())
}

func (g *Source) Project() string {
	return fmt.Sprintf("%s/%s", g.org, g.repo)
}

func (g *Source) Fetch(ctx context.Context, thread *starlark.Thread) (*plugins.SourceResult, error) {
	items := []RepoItem{}
	issues, err := getIssuesForRepo(ctx, g.client, g.org, g.repo, g.issueListOptions)
	if err != nil {
		return nil, fmt.Errorf("fetching issues for GitHub repository %s/%s : %w", g.org, g.repo, err)
	}
	items = append(items, issues...)

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

func newGithubClient() *github.Client {
	token := os.Getenv("SYNKR_GITHUB_TOKEN")
	if token == "" {
		return github.NewClient(http.DefaultClient)
	}

	return github.NewClient(http.DefaultClient).WithAuthToken(token)
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
	ID        int64        `json:"id"`
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
	Milestone string       `json:"milestone"`

	// Only populated on PullRequests
	RequestedReviewers []string `json:"requestedReviewers"`
	Draft              bool     `json:"draft"`
}

func (ri RepoItem) Identifier() string {
	return fmt.Sprintf("%d", ri.ID)
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
	_ = dict.SetKey(starlark.String("milestone"), starlark.String(item.Milestone))

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

func getIssuesForRepo(ctx context.Context, client *github.Client, org, repo string, listOpts *github.IssueListByRepoOptions) ([]RepoItem, error) {
	issues, _, err := client.Issues.ListByRepo(ctx, org, repo, listOpts)
	if err != nil {
		return nil, fmt.Errorf("fetching issues: %v", err)
	}

	out := []RepoItem{}
	for _, issue := range issues {
		if issue == nil {
			continue
		}

		item := RepoItem{
			Type: RepoItemTypeIssue,
		}

		if issue.IsPullRequest() {
			item.Type = RepoItemTypePullRequest

			if issue.PullRequestLinks.URL != nil {
				splits := strings.Split(*issue.PullRequestLinks.URL, "/")
				numStr := splits[len(splits)-1]
				// ignore any errors that happen here because worst case we just get lower
				// fidelity information.
				prNum, err := strconv.Atoi(numStr)
				if err == nil {
					pr, _, _ := client.PullRequests.Get(ctx, org, repo, prNum)
					if pr != nil {
						for _, user := range pr.RequestedReviewers {
							item.RequestedReviewers = append(item.RequestedReviewers, *user.Login)
						}

						item.Draft = *pr.Draft
					}
				}
			}
		}

		if issue.ID != nil {
			item.ID = *issue.ID
		}

		if issue.HTMLURL != nil {
			item.URL = *issue.HTMLURL
		}

		if issue.User != nil && issue.User.Login != nil {
			item.Author = *issue.User.Login
		}

		if issue.Title != nil {
			item.Title = *issue.Title
		}

		if issue.Body != nil {
			item.Body = *issue.Body
		}

		if issue.State != nil {
			item.State = *issue.State
		}

		if issue.Assignees != nil {
			assignees := []string{}
			for _, assignee := range issue.Assignees {
				if assignee == nil {
					continue
				}

				if assignee.Login != nil {
					assignees = append(assignees, *assignee.Login)
				}
			}

			item.Assignees = assignees
		}

		if issue.Labels != nil {
			labels := []string{}
			for _, label := range issue.Labels {
				if label == nil {
					continue
				}

				if label.Name != nil {
					labels = append(labels, *label.Name)
				}
			}

			item.Labels = labels
		}

		if issue.CreatedAt != nil {
			item.Created = issue.CreatedAt.String()
		}

		if issue.UpdatedAt != nil {
			item.Updated = issue.UpdatedAt.String()
		}

		if issue.Comments != nil {
			item.Comments = *issue.Comments
		}

		if issue.Milestone != nil && issue.Milestone.Title != nil {
			item.Milestone = *issue.Milestone.Title
		}

		out = append(out, item)
	}

	return out, nil
}
