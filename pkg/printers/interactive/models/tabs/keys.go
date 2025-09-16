package tabs

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Next key.Binding
	Prev key.Binding
}

var DefaultHorizontalKeyMap = KeyMap{
	Next: key.NewBinding(key.WithKeys("L")),
	Prev: key.NewBinding(key.WithKeys("H")),
}

var DefaultVerticalKeyMap = KeyMap{
	Next: key.NewBinding(key.WithKeys("J")),
	Prev: key.NewBinding(key.WithKeys("K")),
}
