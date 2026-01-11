package tui

import "github.com/charmbracelet/bubbles/key"

var (
	unsubscribe = key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "unsubscribe"),
	)

	deleteAll = key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "delete"),
	)

	archiveAll = key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "archive"),
	)

	trashAll = key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "trash"),
	)

	markRead = key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "mark as read"),
	)

	markSpam = key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "mark as spam"),
	)

	inspect = key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "inspect"),
	)

	goBack = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	)

	toggleHelpMenu = key.NewBinding(
		key.WithKeys("H"),
		key.WithHelp("H", "toggle help"),
	)
)
