package pager

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/v2"
)

type PageChangeFunc func(int)

type Model struct {
	page       int
	totalPages int
	style      lipgloss.Style
	keys       KeyMap
	wrap       bool
	onChange   PageChangeFunc
}

func New(pages int, onChange PageChangeFunc, opts ...Option) *Model {
	m := &Model{
		page:       0,
		totalPages: pages,
		style:      DefaultStyle,
		keys:       DefaultKeyMap,
		onChange:   onChange,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (p *Model) Init() tea.Cmd {
	return nil
}

func (p *Model) Update(m tea.Msg) (*Model, tea.Cmd) {
	switch msg := m.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, p.keys.NextPage) {
			p.increment()
			p.onChange(p.page)
		}

		if key.Matches(msg, p.keys.PrevPage) {
			p.decrement()
			p.onChange(p.page)
		}
	}
	return p, nil
}

func (p *Model) View() string {
	return p.style.Render(fmt.Sprintf("%d / %d", p.page+1, p.totalPages))
}

func (p *Model) increment() {
	if p.page < p.totalPages-1 {
		p.page++
		return
	}

	if p.page == p.totalPages-1 && p.wrap {
		p.page = 0
	}
}

func (p *Model) decrement() {
	if p.page > 0 {
		p.page--
		return
	}

	if p.page == 0 && p.wrap {
		p.page = p.totalPages - 1
	}
}

func (p *Model) Page() int {
	return p.page
}
