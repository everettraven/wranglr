package pager

import "github.com/charmbracelet/lipgloss/v2"

type Option func(p *Model)

func WithWrapping(wrap bool) Option {
	return func(p *Model) {
		p.wrap = wrap
	}
}

func WithStyle(style lipgloss.Style) Option {
	return func(p *Model) {
		p.style = style
	}
}

func WithKeyMap(keyMap KeyMap) Option {
	return func(p *Model) {
		p.keys = keyMap
	}
}
