package jira

import (
	"errors"
	"fmt"
	"time"

	gojira "github.com/andygrunwald/go-jira"
	"github.com/everettraven/wranglr/pkg/modules"
	"go.starlark.net/starlark"
)

type Module struct{}

func (m *Module) String() string        { return "jira" }
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

type searchResult struct {
	Issues []gojira.Issue `json:"issues,omitempty"`
}

func SearchBuiltin() modules.BuiltinFunc {
	return func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var host starlark.String
		var query starlark.String
		var group starlark.String

		err := starlark.UnpackArgs(SearchAttr, args, kwargs,
			"host", &host,
			"query", &query,
			"group?", group,
		)
		if err != nil {
			return nil, err
		}

		client := NewClient(host.GoString())

		issues, err := client.Issues(query.GoString())
		if err != nil {
			return nil, err
		}

		return issuesToStarlark(host.GoString(), group.GoString(), issues...), nil
	}
}

type Item struct {
	issue    gojira.Issue
	status   string
	priority int64
	url      string
	group    string
}

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
	return i.url
}

func (i *Item) Issue() gojira.Issue {
	return i.issue
}

func (i *Item) Group() string {
	if i.group == "" {
		return "Unknown"
	}

	return i.group
}

func (i *Item) String() string        { return "todo" }
func (i *Item) Type() string          { return fmt.Sprintf("%T", i) }
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
	case "assignee":
		if i.issue.Fields.Assignee != nil {
			return starlark.String(i.issue.Fields.Assignee.Name), nil
		}
		return starlark.None, nil
	case "creator":
		if i.issue.Fields.Creator != nil {
			return starlark.String(i.issue.Fields.Creator.Name), nil
		}
		return starlark.None, nil
	case "reporter":
		if i.issue.Fields.Reporter != nil {
			return starlark.String(i.issue.Fields.Reporter.Name), nil
		}
		return starlark.None, nil
	case "type":
		return starlark.String(i.issue.Fields.Type.Name), nil
	case "project":
		return starlark.String(i.issue.Fields.Project.Name), nil
	case "resolution":
		if i.issue.Fields.Resolution != nil {
			return starlark.String(i.issue.Fields.Resolution.Name), nil
		}
		return starlark.None, nil
	case "ticket_priority":
		if i.issue.Fields.Priority != nil {
			return starlark.String(i.issue.Fields.Priority.Name), nil
		}
		return starlark.None, nil
	case "resolution_date":
		return starlark.String(time.Time(i.issue.Fields.Resolutiondate).String()), nil
	case "created":
		return starlark.String(time.Time(i.issue.Fields.Created).String()), nil
	case "due_date":
		return starlark.String(time.Time(i.issue.Fields.Duedate).String()), nil
	case "updated":
		return starlark.String(time.Time(i.issue.Fields.Updated).String()), nil
	case "description":
		return starlark.String(i.issue.Fields.Description), nil
	case "summary":
		return starlark.String(i.issue.Fields.Summary), nil
	case "components":
		elems := []starlark.Value{}
		for _, component := range i.issue.Fields.Components {
			if component == nil {
				continue
			}
			elems = append(elems, starlark.String(component.Name))
		}
		return starlark.NewList(elems), nil
	case "ticket_status":
		if i.issue.Fields.Status != nil {
			return starlark.String(i.issue.Fields.Status.Name), nil
		}
		return starlark.None, nil
	case "fix_versions":
		elems := []starlark.Value{}
		for _, fixVersion := range i.issue.Fields.FixVersions {
			if fixVersion == nil {
				continue
			}
			elems = append(elems, starlark.String(fixVersion.Name))
		}
		return starlark.NewList(elems), nil
	case "affects_versions":
		elems := []starlark.Value{}
		for _, affectsVersion := range i.issue.Fields.AffectsVersions {
			if affectsVersion == nil {
				continue
			}
			elems = append(elems, starlark.String(affectsVersion.Name))
		}
		return starlark.NewList(elems), nil
	case "labels":
		elems := []starlark.Value{}
		for _, label := range i.issue.Fields.Labels {
			elems = append(elems, starlark.String(label))
		}
		return starlark.NewList(elems), nil
	case "epic":
		if i.issue.Fields.Epic != nil {
			return starlark.String(i.issue.Fields.Epic.Name), nil
		}
		return starlark.None, nil
	case "sprint":
		if i.issue.Fields.Sprint != nil {
			return starlark.String(i.issue.Fields.Sprint.Name), nil
		}
		return starlark.None, nil
	default:
		return nil, fmt.Errorf("unknown attribute %q", name)
	}
}

func (i *Item) AttrNames() []string {
	return []string{
		"status",
		"priority",
		"assignee",
		"creator",
		"reporter",
		"type",
		"project",
		"resolution",
		"ticket_priority",
		"resolution_date",
		"created",
		"due_date",
		"updated",
		"description",
		"summary",
		"components",
		"ticket_status",
		"fix_versions",
		"affects_versions",
		"labels",
		"epic",
		"sprint",
	}
}

func (i *Item) SetField(name string, val starlark.Value) error {
	switch name {
	case "status":
		i.status = val.String()
		return nil
	case "priority":
		intType, ok := val.(*starlark.Int)
		if !ok {
			return fmt.Errorf("priority must be an integer but was attempted to be set to type %q", val.Type())
		}
		i64, ok := intType.Int64()
		if !ok {
			return errors.New("priority must be a valid int64, but was not")
		}
		i.priority = i64
		return nil
	default:
		return fmt.Errorf("cannot set field %q", name)
	}
}

func issuesToStarlark(host, group string, issues ...gojira.Issue) starlark.Value {
	elems := []starlark.Value{}
	for _, issue := range issues {
		elems = append(elems, &Item{
			issue: issue,
			url:   fmt.Sprintf("%s/browse/%s", host, issue.Key),
			group: group,
		})
	}

	return starlark.NewList(elems)
}
