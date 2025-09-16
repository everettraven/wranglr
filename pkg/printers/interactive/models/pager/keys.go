package pager

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	NextPage key.Binding
	PrevPage key.Binding
}

var DefaultKeyMap KeyMap = KeyMap{
	NextPage: key.NewBinding(key.WithKeys("l")),
	PrevPage: key.NewBinding(key.WithKeys("h")),
}
