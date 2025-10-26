package interactables

import "github.com/charmbracelet/lipgloss/v2"

var (
	projectStyle     = lipgloss.NewStyle().Foreground(lipgloss.White).Faint(true).Italic(true)
	titleStyle       = lipgloss.NewStyle().Foreground(lipgloss.White).Bold(true)
	bodyStyle        = lipgloss.NewStyle().Foreground(lipgloss.White).Faint(true)
	labelStyle       = lipgloss.NewStyle().Foreground(lipgloss.Black).Align(lipgloss.Center).Padding(0, 1, 0, 1)
	stateOpenStyle   = lipgloss.NewStyle().Foreground(lipgloss.Green)
	stateClosedStyle = lipgloss.NewStyle().Foreground(lipgloss.Red)
	stateMergedStyle = lipgloss.NewStyle().Foreground(lipgloss.Magenta)
)
