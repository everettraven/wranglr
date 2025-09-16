package interactables

import (
	"fmt"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/everettraven/wranglr/pkg/modules/github"
	"github.com/everettraven/wranglr/pkg/printers/interactive/interactables/linkopener"
)

const (
	issueOpen   = ""
	issueClosed = ""
	prOpen      = "󰓂"
	prClosed    = ""
	prMerged    = ""
)

type GitHub struct {
	item *github.Item
}

func NewGitHub(item *github.Item) *GitHub {
	return &GitHub{
		item: item,
	}
}

func (g *GitHub) Priority() int64 {
	return g.item.Priority()
}

func (g *GitHub) Status() string {
	return g.item.Status()
}

func (g *GitHub) Group() string {
	return g.item.Group()
}

func (g *GitHub) Render(width int) string {
	var out strings.Builder

	issue := g.item.Issue()

	symbol := ""
	switch g.item.Type() {
	case string(github.ItemTypeIssue):
		switch issue.State() {
		case "open":
			symbol = stateOpenStyle.Render(issueOpen)
		case "closed":
			symbol = stateClosedStyle.Render(issueClosed)
		}
	case string(github.ItemTypePullRequest):
		switch issue.State() {
		case "open":
			symbol = stateOpenStyle.Render(prOpen)
		case "closed":
			symbol = stateClosedStyle.Render(prClosed)
		case "merged":
			symbol = stateMergedStyle.Render(prMerged)
		}
	}

	prefixRegex := regexp.MustCompile("^https://api.+/repos/")
	prefix := prefixRegex.FindString(issue.RepositoryURL)
	project := strings.TrimPrefix(issue.RepositoryURL, prefix)
	out.WriteString(projectStyle.Render(fmt.Sprintf("%s  %s", "", project)) + "\n\n")

	out.WriteString(titleStyle.Width(width).Render(fmt.Sprintf("%s  %s", symbol, issue.Title)) + "\n\n")

	out.WriteString(fmt.Sprintf(
		"%s %s",
		projectStyle.Render("by"),
		titleStyle.Render(fmt.Sprintf("@%s", issue.Author.Login)),
	))
	out.WriteString("\n\n")

	out.WriteString(" ")
	if len(issue.Assignees) > 0 {
		for _, assignee := range issue.Assignees {
			out.WriteString(titleStyle.Render(fmt.Sprintf("@%s ", assignee.Login)))
		}
	} else {
		out.WriteString(projectStyle.Render("unassigned"))
	}

	out.WriteString("\n\n")

	labelsStr := ""
	for _, label := range issue.Labels {
		labelsStr += labelStyle.Background(lipgloss.Color(fmt.Sprintf("#%s", label.Color))).Render(label.Name) + " "
	}

	if len(labelsStr) > 0 {
		out.WriteString(lipgloss.NewStyle().Width(width).Render(labelsStr))
		out.WriteString("\n")
	}

	bodyOut, _ := glamour.Render(issue.Body, "dark")
	out.WriteString(bodyOut)

	return lipgloss.NewStyle().MarginLeft(2).Render(out.String())
}

func (g *GitHub) Open() tea.Cmd {
	cmd := linkopener.New(g.item.URL()).Open()
	return tea.ExecProcess(cmd, nil)
}
