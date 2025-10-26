package interactables

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/everettraven/wranglr/pkg/modules/jira"
	"github.com/everettraven/wranglr/pkg/printers/interactive/interactables/linkopener"
)

type Jira struct {
	item *jira.Item
}

func NewJira(item *jira.Item) *Jira {
	return &Jira{
		item: item,
	}
}

func (j *Jira) Priority() int64 {
	return j.item.Priority()
}

func (j *Jira) Status() string {
	return j.item.Status()
}

func (j *Jira) Group() string {
	return j.item.Group()
}

func (j *Jira) Render(width int) string {
	var out strings.Builder

	issue := j.item.Issue()

	out.WriteString(projectStyle.Render(fmt.Sprintf("%s %s", "", issue.Key)) + "\n")
	out.WriteString(titleStyle.Width(width).Render(fmt.Sprintf("[%s] %s", issue.Fields.Type.Name, issue.Fields.Summary)) + "\n")

	out.WriteString(fmt.Sprintf(
		"%s %s",
		projectStyle.Render("by"),
		titleStyle.Render(issue.Fields.Reporter.DisplayName),
	))
	out.WriteString("\n")

	out.WriteString("  ")
	if issue.Fields.Assignee != nil {
		out.WriteString(titleStyle.Render(issue.Fields.Assignee.DisplayName))
	} else {
		out.WriteString(projectStyle.Render("unassigned"))
	}
	out.WriteString("\n\n")

	out.WriteString(issue.Fields.Priority.Name + "\n")

	out.WriteString("  ")
	for _, component := range issue.Fields.Components {
		out.WriteString(titleStyle.Render(component.Name + " "))
	}
	out.WriteString("\n\n")

	labelsStr := ""
	for _, label := range issue.Fields.Labels {
		labelsStr += labelStyle.Background(lipgloss.Cyan).Render(label) + " "
	}
	if len(labelsStr) > 0 {
		out.WriteString(lipgloss.NewStyle().Width(width).Render(labelsStr))
		out.WriteString("\n")
	}

	// TODO: Improved rendering for Jira ticket bodies.
	// This will likely include needing to build an
	// Atlassian Document Format parsing library
	// that we can then feed into a glamour-like library
	// for rendering different types of document nodes (headings, codeblocks, etc.).
	bodyOut := issue.Fields.Description
	wrapped := lipgloss.NewStyle().Width(width)
	out.WriteString(wrapped.Render(bodyStyle.Render(bodyOut)))

	return lipgloss.NewStyle().MarginLeft(2).Render(out.String())
}

func (j *Jira) Open() tea.Cmd {
	cmd := linkopener.New(j.item.URL()).Open()
	return tea.ExecProcess(cmd, nil)
}
