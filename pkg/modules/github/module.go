package github

import (
	"context"
	"errors"
	"fmt"

	"github.com/cli/cli/v2/pkg/search"
	"go.starlark.net/starlark"

	"github.com/everettraven/wranglr/pkg/modules"
)

// TODO:
// - Need to figure out how to make this information available to the
// Starlark LSP.

type Module struct{}

func (m *Module) String() string        { return "github" }
func (m *Module) Type() string          { return "Module" }
func (m *Module) Truth() starlark.Bool  { return starlark.False }
func (m *Module) Freeze()               {}
func (m *Module) Hash() (uint32, error) { return 0, fmt.Errorf("hashing not yet implemented") }

const SearchAttr = "search"

func (m *Module) Attr(name string) (starlark.Value, error) {
	switch name {
	case SearchAttr:
		return starlark.NewBuiltin(SearchAttr, SearchBuiltin()), nil
	default:
		return nil, fmt.Errorf("unknown attribute %q", name)
	}
}

func (m *Module) AttrNames() []string {
	return []string{
		SearchAttr,
	}
}

func SearchBuiltin() modules.BuiltinFunc {
	return func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var host starlark.String
		var query starlark.String
		var group starlark.String

		err := starlark.UnpackArgs(SearchAttr, args, kwargs,
			"host?", &host,
			"query", &query,
			"group?", &group,
		)
		if err != nil {
			return nil, err
		}

		hostValue := "github.com"
		if host.GoString() != "" {
			hostValue = host.GoString()
		}

		ghClient := NewClient(hostValue)

		issues, err := ghClient.Issues(context.TODO(), query.GoString())
		if err != nil {
			return nil, err
		}

		return issuesToStarlark(group.GoString(), issues...), nil
	}
}

type Item struct {
	issue    search.Issue
	status   string
	priority int64
	group    string
}

type ItemType string

const (
	ItemTypeIssue       ItemType = "issue"
	ItemTypePullRequest ItemType = "pullrequest"
)

func (i *Item) Priority() int64 {
	return i.priority
}

func (i *Item) Status() string {
	if i.status == "" {
		return "Unknown"
	}

	return i.status
}

func (i *Item) URL() string {
	return i.issue.URL
}

func (i *Item) Group() string {
	if i.group == "" {
		return "Unknown"
	}
	return i.group
}

func (i *Item) Issue() search.Issue {
	return i.issue
}

func (i *Item) String() string { return "todo" }
func (i *Item) Type() string {
	if i.issue.PullRequest.URL != "" {
		return string(ItemTypePullRequest)
	}

	return string(ItemTypeIssue)
}
func (i *Item) Truth() starlark.Bool  { return starlark.False }
func (i *Item) Freeze()               {}
func (i *Item) Hash() (uint32, error) { return 0, fmt.Errorf("hashing not yet implemented") }

func (i *Item) Attr(name string) (starlark.Value, error) {
	switch name {
	case "status":
		return starlark.String(i.Status()), nil
	case "priority":
		return starlark.MakeInt64(i.Priority()), nil
	case "group":
		return starlark.String(i.Group()), nil
	case "assignees":
		elems := []starlark.Value{}
		for _, assignee := range i.issue.Assignees {
			elems = append(elems, starlark.String(assignee.Login))
		}
		return starlark.NewList(elems), nil
	case "author":
		return starlark.String(i.issue.Author.Login), nil
	case "author_association":
		return starlark.String(i.issue.AuthorAssociation), nil
	case "body":
		return starlark.String(i.issue.Body), nil
	case "closed_at":
		return starlark.String(i.issue.ClosedAt.String()), nil
	case "comments":
		return starlark.MakeInt(i.issue.CommentsCount), nil
	case "created_at":
		return starlark.String(i.issue.CreatedAt.String()), nil
	case "labels":
		elems := []starlark.Value{}
		for _, label := range i.issue.Labels {
			elems = append(elems, starlark.String(label.Name))
		}
		return starlark.NewList(elems), nil
	case "locked":
		return starlark.Bool(i.issue.IsLocked), nil
	case "number":
		return starlark.MakeInt(i.issue.Number), nil
	case "pull_request":
		// TODO: make this return a custom starlark.Value that implements attributes.
		// For now this will just be a dict.
		dict := starlark.NewDict(2)
		err := dict.SetKey(starlark.String("url"), starlark.String(i.issue.PullRequest.URL))
		if err != nil {
			return starlark.None, err
		}

		err = dict.SetKey(starlark.String("merged_at"), starlark.String(i.issue.PullRequest.MergedAt.String()))
		if err != nil {
			return starlark.None, err
		}

		return dict, nil
	case "state":
		return starlark.String(i.issue.State()), nil
	case "state_reason":
		return starlark.String(i.issue.StateReason), nil
	case "title":
		return starlark.String(i.issue.Title), nil
	case "updated_at":
		return starlark.String(i.issue.UpdatedAt.String()), nil
	default:
		return nil, fmt.Errorf("unknown attribute %q", name)
	}
}

func (i *Item) AttrNames() []string {
	return []string{
		"status",
		"priority",
		"group",
		"assignees",
		"author",
		"author_association",
		"body",
		"closed_at",
		"comments",
		"created_at",
		"labels",
		"locked",
		"number",
		"pull_request",
		"state",
		"state_reason",
		"title",
		"updated_at",
	}
}

func (i *Item) SetField(name string, val starlark.Value) error {
	switch name {
	case "status":
		i.status = val.String()
		return nil
	case "priority":
		intType, ok := val.(starlark.Int)
		if !ok {
			return fmt.Errorf("priority must be an integer but was attempted to be set to type %q", val.Type())
		}
		i64, ok := intType.Int64()
		if !ok {
			return errors.New("priority must be a valid int64, but was not")
		}
		i.priority = i64
		return nil
	case "group":
		i.group = val.String()
		return nil
	default:
		return fmt.Errorf("cannot set field %q", name)
	}
}

func issuesToStarlark(group string, issues ...search.Issue) starlark.Value {
	elems := []starlark.Value{}
	for _, issue := range issues {
		elems = append(elems, &Item{
			issue: issue,
			group: group,
		})
	}

	return starlark.NewList(elems)
}
