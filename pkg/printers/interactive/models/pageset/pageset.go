package pageset

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/everettraven/wranglr/pkg/printers/interactive/models/pager"
)

type Renderable interface {
	Render(width int) string
}

type Openable interface {
	Open() tea.Cmd
}

type Page interface {
	Renderable
	Openable
}

var DefaultStyle = lipgloss.NewStyle().Margin(0, 0, 1, 2)

type PageSet struct {
	viewportModel viewport.Model
	pager         *pager.Model
	pages         []Page
	style         lipgloss.Style
}

// TODO: optionality
func New(pages ...Page) *PageSet {
	ps := &PageSet{
		viewportModel: viewport.New(100, 100),
		pages:         pages,
		style:         DefaultStyle,
	}

	ps.pager = pager.New(
		len(pages),
		func(i int) {
			ps.viewportModel.SetContent(ps.pages[i].Render(ps.viewportModel.Width))
		},
	)

	return ps
}

func (ps *PageSet) Init() tea.Cmd {
	ps.viewportModel.SetContent(ps.pages[ps.pager.Page()].Render(ps.viewportModel.Width))
	return nil
}

func (ps *PageSet) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		ps.viewportModel.Width = msg.Width
		ps.viewportModel.Height = msg.Height - lipgloss.Height(ps.pager.View())
		return ps, nil

	case tea.KeyMsg:
		switch msg.String() {
		// TODO: make this "action-oriented"
		// where a page can define "actions" that it
		// supports and then can be dynamically supported.
		case "o":
			return ps, ps.pages[ps.pager.Page()].Open()
		}
	}

	var cmd tea.Cmd

	ps.viewportModel, _ = ps.viewportModel.Update(message)
	ps.pager, cmd = ps.pager.Update(message)

	return ps, cmd
}

func (ps *PageSet) View() string {
	return ps.style.Render(
		lipgloss.JoinVertical(
			lipgloss.Top,
			ps.viewportModel.View(),
			ps.pager.View(),
		),
	)
}
