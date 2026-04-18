package ui

import "github.com/charmbracelet/bubbles/key"

// keyMap defines all key bindings used across states.
type keyMap struct {
	Submit key.Binding
	Quit   key.Binding
	Retry  key.Binding
	Skip   key.Binding
}

var keys = keyMap{
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "submit"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q/esc", "quit"),
	),
	Retry: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "retry interactively"),
	),
	Skip: key.NewBinding(
		key.WithKeys("s", "enter"),
		key.WithHelp("s/enter", "skip"),
	),
}
