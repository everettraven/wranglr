package interactive

import (
	"cmp"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/everettraven/wranglr/pkg/printers/interactive/models/pageset"
	"github.com/everettraven/wranglr/pkg/printers/interactive/models/tabs"
)

type Prioritizable interface {
	Priority() int64
}

type Statuser interface {
	Status() string
}

type Renderable interface {
	Render(width int) string
}

type Openable interface {
	Open() tea.Cmd
}

type Grouper interface {
	Group() string
}

type Interactable interface {
	Prioritizable
	Statuser
	Renderable
	Openable
	Grouper
}

func NewRoot(entries ...Interactable) *Root {
	// sort by priority score
	slices.SortFunc(entries, func(a, b Interactable) int {
		return cmp.Compare(a.Priority(), b.Priority())
	})

	// to sort in inverse order where higher priority is first
	slices.Reverse(entries)

	groupedStatusedPages := pagesByGroupAndStatus(entries...)

	return &Root{
		tabs: groupedTabs(groupedStatusedPages),
	}
}

func groupedTabs(groups GroupedStatusedPages) *tabs.Model {
	tabList := []tabs.Tab{}
	for k, v := range groups {
		normalizedKey := strings.TrimSuffix(strings.TrimPrefix(k, "\""), "\"")
		tabList = append(tabList, tabs.Tab{
			Name:  normalizedKey,
			Model: statusTabs(v),
		})
	}

	// Tabs sorted for determinism
	slices.SortFunc(tabList, func(a, b tabs.Tab) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return tabs.New(
		tabList,
		tabs.WithWrapping(true),
		tabs.WithDisplayFormat(tabs.DisplayFormatVertical),
		tabs.WithVerticalWidth(15),
	)
}

func statusTabs(statuses StatusedPages) *tabs.Model {
	tabList := []tabs.Tab{}
	for k, v := range statuses {
		// Don't render empty lists of data
		if len(v) == 0 {
			continue
		}

		normalizedKey := strings.TrimSuffix(strings.TrimPrefix(k, "\""), "\"")
		tabList = append(tabList, tabs.Tab{
			Name:  normalizedKey,
			Model: pageset.New(v...),
		})
	}

	// Tabs sorted for determinism
	slices.SortFunc(tabList, func(a, b tabs.Tab) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return tabs.New(
		tabList,
		tabs.WithWrapping(true),
	)
}

type GroupedStatusedPages map[string]StatusedPages

type StatusedPages map[string][]pageset.Page

func pagesByGroupAndStatus(entries ...Interactable) GroupedStatusedPages {
	grouped := map[string][]Interactable{}
	for _, entry := range entries {
		grouped[entry.Group()] = append(grouped[entry.Group()], entry)
	}

	ultraGrouped := make(GroupedStatusedPages)
	for k, v := range grouped {
		ultraGrouped[k] = pagesByStatus(v...)
	}

	return ultraGrouped
}

func pagesByStatus(entries ...Interactable) StatusedPages {
	slicesByStatus := map[string][]pageset.Page{}

	for _, entry := range entries {
		slicesByStatus[entry.Status()] = append(slicesByStatus[entry.Status()], entry)
	}

	return slicesByStatus
}

type Root struct {
	tabs *tabs.Model
}

func (r *Root) Init() tea.Cmd {
	return r.tabs.Init()
}

func (r *Root) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return r, tea.Quit
		}
	}

	var cmd tea.Cmd
	var tabsModel tea.Model

	tabsModel, cmd = r.tabs.Update(message)

	r.tabs = tabsModel.(*tabs.Model)
	return r, cmd
}

func (r *Root) View() string {
	return r.tabs.View()
}
