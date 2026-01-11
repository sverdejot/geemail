package tui

import "github.com/charmbracelet/bubbles/key"

var (
	unsubscribe = key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "unsubscribe from list"),
	)

	deleteAll = key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "delete all mails from this sender"),
	)

	archiveAll = key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "archive all mails from this sender"),
	)

	trashAll = key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "trash all mails from this sender"),
	)

	markRead = key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "mark all mails as read"),
	)

	markSpam = key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "mark all mails as spam"),
	)

	toggleHelpMenu = key.NewBinding(
		key.WithKeys("H"),
		key.WithHelp("H", "toggle help"),
	)
)
