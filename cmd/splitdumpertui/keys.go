package main

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Incr    key.Binding
	Decr    key.Binding
	Execute key.Binding
	Quit    key.Binding
}

var DefaultKeyMap = KeyMap{
	Incr:    key.NewBinding(key.WithKeys("up"), key.WithHelp("↑", "")),
	Decr:    key.NewBinding(key.WithKeys("down"), key.WithHelp("↓", "")),
	Execute: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "")),
	Quit:    key.NewBinding(key.WithKeys("esc", "ctrl+c"), key.WithHelp("ctrl+c", "")),
}
