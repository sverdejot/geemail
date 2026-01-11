package tui

import (
	"github.com/sverdejot/geemail/internal/inbox"
)

// Mail streaming messages - emitted during mail loading
type mailReceivedMsg struct {
	mail inbox.RawMail
}

type mailStreamReadyMsg struct {
	stream <-chan inbox.RawMail
}

type mailStreamCompleteMsg struct{}

type mailStreamErrorMsg struct {
	err error
}

// Intent messages - emitted by mailList when user takes action
type unsubscribeRequestMsg struct {
	mail inbox.MailingList
	idx  int
}

type deleteRequestMsg struct {
	mail inbox.MailingList
	idx  int
}

type archiveRequestMsg struct {
	mail inbox.MailingList
	idx  int
}

type trashRequestMsg struct {
	mail inbox.MailingList
	idx  int
}

// Result messages - emitted after async operations complete
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

type archiveCompleteMsg struct {
	mail inbox.MailingList
	idx  int
	err  error
}

type trashCompleteMsg struct {
	mail inbox.MailingList
	idx  int
	err  error
}

// Status message for user feedback
type statusMsg struct {
	text string
}
