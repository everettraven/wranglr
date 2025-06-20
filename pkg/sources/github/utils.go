package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v71/github"
	"go.starlark.net/starlark"
)

func repoItemToStarlarkDict(item RepoItem) *starlark.Dict {
	dict := &starlark.Dict{}

	dict.SetKey(starlark.String("id"), starlark.MakeInt64(item.ID))
	dict.SetKey(starlark.String("url"), starlark.String(item.URL))
	dict.SetKey(starlark.String("author"), starlark.String(item.Author))
	dict.SetKey(starlark.String("type"), starlark.String(item.Type))
	dict.SetKey(starlark.String("title"), starlark.String(item.Title))
	dict.SetKey(starlark.String("body"), starlark.String(item.Body))
	dict.SetKey(starlark.String("state"), starlark.String(item.State))
	dict.SetKey(starlark.String("labels"), starlark.NewList(stringArrayToStarlarkValueArray(item.Labels...)))
	dict.SetKey(starlark.String("assignees"), starlark.NewList(stringArrayToStarlarkValueArray(item.Assignees...)))

	return dict
}

func stringArrayToStarlarkValueArray(in ...string) []starlark.Value {
	out := []starlark.Value{}
	for _, i := range in {
		out = append(out, starlark.String(i))
	}

	return out
}

// TODO: Probably makes sense long term to have some sort of common function for
// translating issues/prs into the RepoItem format.
// Alternatively, just use the github.Issue format.
//
// TODO: Update this to allow for authenticated and unauthenticated access to the GH API.
// Authenticated would allow access to private repositories, higher rate limits, and enterprise GH instances.
// Unauthenticated would still be the default.
func getPullRequestsForRepo(ctx context.Context, org, repo string) ([]RepoItem, error) {
	client := github.NewClient(nil)
	pullRequests, _, err := client.PullRequests.List(ctx, org, repo, nil)
	if err != nil {
		return nil, fmt.Errorf("fetching pull requests: %v", err)
	}

	out := []RepoItem{}
	for _, issue := range pullRequests {
		if issue == nil {
			continue
		}

		item := RepoItem{
			Type: RepoItemTypePullRequest,
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

		out = append(out, item)
	}

	return out, nil
}

func getIssuesForRepo(ctx context.Context, org, repo string) ([]RepoItem, error) {
	client := github.NewClient(nil)
	issues, _, err := client.Issues.ListByRepo(ctx, org, repo, nil)
	if err != nil {
		return nil, fmt.Errorf("fetching issues: %v", err)
	}

	out := []RepoItem{}
	for _, issue := range issues {
		if issue == nil {
			continue
		}

		// ignore PRs - we fetch those separately for now
		if issue.IsPullRequest() {
			continue
		}

		item := RepoItem{
			Type: RepoItemTypeIssue,
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

		out = append(out, item)
	}

	return out, nil
}
