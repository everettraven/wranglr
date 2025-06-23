package github

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/go-github/v71/github"
	"go.starlark.net/starlark"
)

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
	_ = dict.SetKey(starlark.String("mentions"), starlark.NewList(stringArrayToStarlarkValueArray(item.Mentions...)))

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

// TODO: Probably makes sense long term to have some sort of common function for
// translating issues/prs into the RepoItem format.
// Alternatively, just use the github.Issue format.
func getPullRequestsForRepo(ctx context.Context, client *github.Client, org, repo string, includeMentions bool) ([]RepoItem, error) {
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

		if issue.RequestedReviewers != nil {
			reviewers := []string{}
			for _, reviewer := range issue.RequestedReviewers {
				if reviewer == nil {
					continue
				}

				if reviewer.Login != nil {
					reviewers = append(reviewers, *reviewer.Login)
				}
			}

			item.RequestedReviewers = reviewers
		}

		if issue.Draft != nil {
			item.Draft = *issue.Draft
		}

		if includeMentions {
			mentions, err := getMentionsForItemID(ctx, client, org, repo, *issue.Number)
			if err != nil {
				return nil, fmt.Errorf("getting mentions for pull request: %w", err)
			}

			item.Mentions = mentions
		}

		out = append(out, item)
	}

	return out, nil
}

func getIssuesForRepo(ctx context.Context, client *github.Client, org, repo string, includeMentions bool) ([]RepoItem, error) {
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

		if includeMentions {
			mentions, err := getMentionsForItemID(ctx, client, org, repo, *issue.Number)
			if err != nil {
				return nil, fmt.Errorf("getting mentions for pull request: %w", err)
			}

			item.Mentions = mentions
		}

		out = append(out, item)
	}

	return out, nil
}

func getMentionsForItemID(ctx context.Context, client *github.Client, org, repo string, id int) ([]string, error) {
	// get comments for item
	comments, _, err := client.Issues.ListComments(ctx, org, repo, id, &github.IssueListCommentsOptions{})
	if err != nil {
		return nil, err
	}

	// get mentions from comments
	return mentionsFromComments(comments...), nil
}

var mentionRegex = regexp.MustCompile(`(@\w+)`)

func mentionsFromComments(comments ...*github.IssueComment) []string {
	mentionsMap := map[string]struct{}{}

	for _, comment := range comments {
		if comment == nil {
			continue
		}

		body := comment.GetBody()
		commentMentions := mentionRegex.FindAllString(body, -1)
		for _, mention := range commentMentions {
			sanitizedMention := strings.TrimPrefix(mention, "@")
			if _, ok := mentionsMap[sanitizedMention]; !ok {
				mentionsMap[sanitizedMention] = struct{}{}
			}
		}

	}

	mentions := []string{}
	for k := range mentionsMap {
		mentions = append(mentions, k)
	}

	return mentions
}
