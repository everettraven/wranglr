package tabs

import "github.com/charmbracelet/lipgloss/v2"

var BaseStyle = lipgloss.NewStyle().Margin(1, 0, 0, 1)

var DefaultHorizontalStyles = Styles{
	Active:   BaseStyle.Bold(true).BorderBottom(true).BorderStyle(lipgloss.NormalBorder()),
	Inactive: BaseStyle.Faint(true),
}

var DefaultVerticalStyles = Styles{
	Active:   BaseStyle.Bold(true).BorderRight(true).BorderStyle(lipgloss.NormalBorder()),
	Inactive: BaseStyle.Faint(true),
}

type Styles struct {
	Active   lipgloss.Style
	Inactive lipgloss.Style
}
