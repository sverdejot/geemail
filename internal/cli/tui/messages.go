package tui

import (
	"github.com/sverdejot/geemail/internal/inbox"
)

// emitted by mailList when user takes action
type unsubscribeRequestMsg struct {
	mail inbox.MailingList
	idx  int
}

type deleteRequestMsg struct {
	mail inbox.MailingList
	idx  int
}

// emitted after async operations complete
type unsubscribeCompleteMsg struct {
	mail inbox.MailingList
	idx  int
	err  error
}

type deleteCompleteMsg struct {
	mail inbox.MailingList
	idx  int
	err  error
}

// Status message for user feedback
type statusMsg struct {
	text string
}
