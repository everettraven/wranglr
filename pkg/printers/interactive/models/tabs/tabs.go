package tabs

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/v2"
)

type Tab struct {
	Name  string
	Model tea.Model
}

type DisplayFormat string

const (
	DisplayFormatVertical   DisplayFormat = "vertical"
	DisplayFormatHorizontal DisplayFormat = "horizontal"
)

type Model struct {
	tabs    []Tab
	styles  Styles
	display DisplayFormat
	keys    KeyMap
	idx     int
	wrap    bool
	width   int
	height  int
}

func New(tabs []Tab, opts ...Options) *Model {
	tabsModel := &Model{
		tabs:    tabs,
		styles:  DefaultHorizontalStyles,
		display: DisplayFormatHorizontal,
		keys:    DefaultHorizontalKeyMap,
	}

	for _, opt := range opts {
		opt(tabsModel)
	}

	return tabsModel
}

func (t *Model) Init() tea.Cmd {
	cmds := []tea.Cmd{}
	for _, tab := range t.tabs {
		cmds = append(cmds, tab.Model.Init())
	}
	return tea.Batch(cmds...)
}

func (t *Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		adjustedWindowSizeMsg := msg
		if t.display == DisplayFormatHorizontal {
			t.width = msg.Width
			adjustedWindowSizeMsg.Height = msg.Height - lipgloss.Height(t.RenderTabs())
		}

		if t.display == DisplayFormatVertical {
			t.height = msg.Height
			adjustedWindowSizeMsg.Width = msg.Width - lipgloss.Width(t.RenderTabs())
		}

		for i := range t.tabs {
			t.tabs[i].Model, _ = t.tabs[i].Model.Update(adjustedWindowSizeMsg)
		}

		return t, nil

	case tea.KeyMsg:
		if key.Matches(msg, t.keys.Next) {
			t.increment()
		}

		if key.Matches(msg, t.keys.Prev) {
			t.decrement()
		}
	}

	var cmd tea.Cmd
	t.tabs[t.idx].Model, cmd = t.tabs[t.idx].Model.Update(message)
	return t, cmd
}

func (t *Model) View() string {
	tabs := ""
	if len(t.tabs) > 1 {
		tabs = t.RenderTabs()
	}

	if t.display == DisplayFormatHorizontal {
		return lipgloss.JoinVertical(lipgloss.Top, tabs, t.tabs[t.idx].Model.View())
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, tabs, t.tabs[t.idx].Model.View())
}

const (
	arrowUp    = ""
	arrowDown  = ""
	arrowLeft  = ""
	arrowRight = ""
)

var horizontalArrowStyle = lipgloss.NewStyle().Faint(true).Margin(1, 0, 0, 1)

func (t *Model) RenderTabs() string {
	if t.display == DisplayFormatHorizontal {
		return t.renderHorizontal()
	}

	return t.renderVertical()
}

func (t *Model) renderHorizontal() string {
	leftArrow := horizontalArrowStyle.Render(arrowLeft)
	rightArrow := horizontalArrowStyle.Render(arrowRight)

	pages := pagesForTabs(
		t.tabs,
		func(s ...string) string {
			return lipgloss.JoinHorizontal(lipgloss.Top, s...)
		},
		func(s string) bool {
			joined := lipgloss.JoinHorizontal(lipgloss.Top, leftArrow, s, rightArrow)
			return lipgloss.Width(joined) > t.width-5
		},
		func(i int, s string) string {
			if i == t.idx {
				return t.styles.Active.Render(s)
			}
			return t.styles.Inactive.Render(s)
		},
	)

	var pageToRender *page
	for _, page := range pages {
		if page.startIndex <= t.idx && page.endIndex >= t.idx {
			pageToRender = page
			break
		}
	}

	tabs := ""
	for i, tab := range t.tabs {
		if i < pageToRender.startIndex || i > pageToRender.endIndex {
			continue
		}

		rendered := t.styles.Inactive.Render(tab.Name)

		if i == t.idx {
			rendered = t.styles.Active.Render(tab.Name)
		}

		tabs = lipgloss.JoinHorizontal(lipgloss.Top, tabs, rendered)
	}

	if len(pages) > 1 {
		return lipgloss.JoinHorizontal(lipgloss.Top, leftArrow, tabs, rightArrow)
	}

	return tabs
}

func (t *Model) renderVertical() string {
	upArrow := horizontalArrowStyle.Render(arrowUp)
	downArrow := horizontalArrowStyle.Render(arrowDown)

	pages := pagesForTabs(
		t.tabs,
		func(s ...string) string {
			return lipgloss.JoinVertical(lipgloss.Top, s...)
		},
		func(s string) bool {
			joined := lipgloss.JoinVertical(lipgloss.Top, upArrow, s, downArrow)
			return lipgloss.Height(joined) > t.height-5
		},
		func(i int, s string) string {
			if i == t.idx {
				return t.styles.Active.Render(s)
			}
			return t.styles.Inactive.Render(s)
		},
	)

	var pageToRender *page
	for _, page := range pages {
		if page.startIndex <= t.idx && page.endIndex >= t.idx {
			pageToRender = page
			break
		}
	}

	tabs := ""
	for i, tab := range t.tabs {
		if i < pageToRender.startIndex || i > pageToRender.endIndex {
			continue
		}

		rendered := t.styles.Inactive.Width(t.width).Render(tab.Name)

		if i == t.idx {
			rendered = t.styles.Active.Width(t.width).Render(tab.Name)
		}

		tabs = lipgloss.JoinVertical(lipgloss.Top, tabs, rendered)
	}

	if len(pages) > 1 {
		return lipgloss.JoinVertical(lipgloss.Top, upArrow, tabs, downArrow)
	}

	return tabs
}

func (t *Model) increment() {
	if t.idx < len(t.tabs)-1 {
		t.idx++
		return
	}

	if t.idx == len(t.tabs)-1 && t.wrap {
		t.idx = 0
	}
}

func (t *Model) decrement() {
	if t.idx > 0 {
		t.idx--
		return
	}

	if t.idx == 0 && t.wrap {
		t.idx = len(t.tabs) - 1
	}
}
