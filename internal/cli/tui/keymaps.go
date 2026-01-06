package tui

import "github.com/charmbracelet/bubbles/key"

var (
	unsubscribe = key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "unsubscribe from list"),
	)

	toggleHelpMenu = key.NewBinding(
		key.WithKeys("H"),
		key.WithHelp("H", "toggle help"),
	)
)
