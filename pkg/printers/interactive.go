package printers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/everettraven/wranglr/pkg/modules/github"
	"github.com/everettraven/wranglr/pkg/modules/jira"
	"github.com/everettraven/wranglr/pkg/printers/interactive"
	"github.com/everettraven/wranglr/pkg/printers/interactive/interactables"
	"go.starlark.net/starlark"
)

type Interactive struct{}

func (i *Interactive) Print(results ...starlark.Value) error {
	interactableResults := []interactive.Interactable{}
	for _, result := range results {
		switch item := result.(type) {
		case *jira.Item:
			interactableResults = append(interactableResults, interactables.NewJira(item))
		case *github.Item:
			interactableResults = append(interactableResults, interactables.NewGitHub(item))
		}
	}
	r := interactive.NewRoot(interactableResults...)

	p := tea.NewProgram(r, tea.WithAltScreen())

	_, err := p.Run()
	if err != nil {
		return err
	}

	return nil
}
