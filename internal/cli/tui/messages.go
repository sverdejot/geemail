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

type markReadRequestMsg struct {
	mail inbox.MailingList
	idx  int
}

type markSpamRequestMsg struct {
	mail inbox.MailingList
	idx  int
}

type inspectRequestMsg struct {
	mail inbox.MailingList
}

type goBackMsg struct{}

// Single mail intent messages - emitted by inspectModel when user takes action
type singleDeleteRequestMsg struct {
	mail inbox.RawMail
	idx  int
}

type singleArchiveRequestMsg struct {
	mail inbox.RawMail
	idx  int
}

type singleTrashRequestMsg struct {
	mail inbox.RawMail
	idx  int
}

type singleMarkReadRequestMsg struct {
	mail inbox.RawMail
	idx  int
}

type singleMarkSpamRequestMsg struct {
	mail inbox.RawMail
	idx  int
}

// Single mail result messages
type singleDeleteCompleteMsg struct {
	mail inbox.RawMail
	idx  int
	err  error
}

type singleArchiveCompleteMsg struct {
	mail inbox.RawMail
	idx  int
	err  error
}

type singleTrashCompleteMsg struct {
	mail inbox.RawMail
	idx  int
	err  error
}

type singleMarkReadCompleteMsg struct {
	mail inbox.RawMail
	idx  int
	err  error
}

type singleMarkSpamCompleteMsg struct {
	mail inbox.RawMail
	idx  int
	err  error
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

type markReadCompleteMsg struct {
	mail inbox.MailingList
	idx  int
	err  error
}

type markSpamCompleteMsg struct {
	mail inbox.MailingList
	idx  int
	err  error
}

// Status message for user feedback
type statusMsg struct {
	text string
}
