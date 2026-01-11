package tui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sverdejot/geemail/internal/gmail"
	"github.com/sverdejot/geemail/internal/inbox"
)

type state int

const (
	loading state = iota
	ready
)

type rootModel struct {
	state               state
	progress            *mailLoadingProgress
	list                mailList
	ctx                 context.Context
	svc                 *gmail.MailService
	mails               []inbox.RawMail
	mailStream          <-chan inbox.RawMail
	width               int
	height              int
	dryRun              bool
	operationInProgress bool
	currentOperation    string
}

func NewRoot(ctx context.Context, svc *gmail.MailService, dryRun bool) (*rootModel, error) {
	total, err := svc.GetTotalUnreads(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get total unread messages: %w", err)
	}

	mails := make([]inbox.RawMail, 0)
	pg := NewProgressModel(total)

	return &rootModel{
		state:    loading,
		progress: pg,
		ctx:      ctx,
		svc:      svc,
		mails:    mails,
		dryRun:   dryRun,
	}, nil
}

func (m *rootModel) Init() tea.Cmd {
	return tea.Batch(
		m.progress.Init(),
		m.startLoading(),
	)
}

func (m *rootModel) startLoading() tea.Cmd {
	return func() tea.Msg {
		stream, err := m.svc.StreamUnreadMessages(m.ctx)
		if err != nil {
			return mailStreamErrorMsg{err: err}
		}

		return mailStreamReadyMsg{stream: stream}
	}
}

func (m *rootModel) readNextMail() tea.Cmd {
	return func() tea.Msg {
		mail, ok := <-m.mailStream
		if !ok {
			return mailStreamCompleteMsg{}
		}

		return mailReceivedMsg{mail: mail}
	}
}

func (m *rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.state == loading {
			_, cmd := m.progress.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			updatedModel, cmd := m.list.Update(msg)
			if updatedList, ok := updatedModel.(mailList); ok {
				m.list = updatedList
			}
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)

	case mailStreamReadyMsg:
		if m.state == loading {
			// Store the stream and start reading from it
			m.mailStream = msg.stream
			return m, m.readNextMail()
		}

	case mailReceivedMsg:
		if m.state == loading {
			m.mails = append(m.mails, msg.mail)

			_, progressCmd := m.progress.Update(msg)

			return m, tea.Batch(progressCmd, m.readNextMail())
		}

	case mailStreamCompleteMsg:
		if m.state == loading {
			// Transform to endMsg for existing logic
			return m, func() tea.Msg { return endMsg{} }
		}

	case mailStreamErrorMsg:
		// Handle stream error (could show error message or quit)
		return m, tea.Quit

	case progressMsg:
		if m.state == loading {
			_, cmd := m.progress.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

	case endMsg:
		if m.state == loading {
			rawMailList := inbox.RawMailList(m.mails)
			mailingLists := inbox.GetMailingList(rawMailList)

			m.list = NewModel(mailingLists)
			if m.width > 0 && m.height > 0 {
				updatedModel, sizeCmd := m.list.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
				if updatedList, ok := updatedModel.(mailList); ok {
					m.list = updatedList
				}
				cmds = append(cmds, sizeCmd)
			}
			m.state = ready

			cmd := m.list.Init()
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

	case unsubscribeRequestMsg:
		if m.operationInProgress {
			return m, m.statusCmd(fmt.Sprintf("Operation '%s' already in progress...", m.currentOperation))
		}

		if m.dryRun {
			return m, m.statusCmd(fmt.Sprintf("[DRY RUN] Would unsubscribe from %s", msg.mail.From))
		}

		m.operationInProgress = true
		m.currentOperation = "unsubscribe"

		return m, func() tea.Msg {
			err := msg.mail.Unsubscribe(m.ctx)
			return unsubscribeCompleteMsg{
				mail: msg.mail,
				idx:  msg.idx,
				err:  err,
			}
		}

	case unsubscribeCompleteMsg:
		m.operationInProgress = false
		m.currentOperation = ""

		if msg.err != nil {
			return m, m.handleUnsubscribeError(msg.mail, msg.err)
		}

		// Unsubscribe succeeded, now delete the emails
		m.operationInProgress = true
		m.currentOperation = "delete"

		return m, func() tea.Msg {
			err := m.svc.BulkDelete(m.ctx, msg.mail.UnreadMessagesIDs)
			return deleteCompleteMsg{
				mail: msg.mail,
				idx:  msg.idx,
				err:  err,
			}
		}

	case deleteRequestMsg:
		if m.operationInProgress {
			return m, m.statusCmd(fmt.Sprintf("Operation '%s' already in progress...", m.currentOperation))
		}

		if m.dryRun {
			return m, m.statusCmd(fmt.Sprintf("[DRY RUN] Would delete (%d) mails from %s", msg.mail.TotalUnreads, msg.mail.From))
		}

		m.operationInProgress = true
		m.currentOperation = "delete"

		return m, func() tea.Msg {
			err := m.svc.BulkDelete(m.ctx, msg.mail.UnreadMessagesIDs)
			return deleteCompleteMsg{
				mail: msg.mail,
				idx:  msg.idx,
				err:  err,
			}
		}

	case deleteCompleteMsg:
		m.operationInProgress = false
		m.currentOperation = ""

		if msg.err != nil {
			return m, m.statusCmd(fmt.Sprintf("Error deleting (%d) mails from %s. Please, try again later.", msg.mail.TotalUnreads, msg.mail.From))
		}

		m.list.list.RemoveItem(msg.idx)
		return m, m.statusCmd(fmt.Sprintf("Deleted (%d) mails from %s", msg.mail.TotalUnreads, msg.mail.From))

	case statusMsg:
		updatedModel, cmd := m.list.Update(m.list.list.NewStatusMessage(msg.text))
		if updatedList, ok := updatedModel.(mailList); ok {
			m.list = updatedList
		}
		return m, cmd

	case tea.KeyMsg:
		if m.state == ready {
			updatedModel, cmd := m.list.Update(msg)
			if updatedList, ok := updatedModel.(mailList); ok {
				m.list = updatedList
			}
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}
		if m.state == loading {
			if msg.String() == "q" || msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
		}

	default:
		if m.state == loading {
			_, cmd := m.progress.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		} else {
			updatedModel, cmd := m.list.Update(msg)
			if updatedList, ok := updatedModel.(mailList); ok {
				m.list = updatedList
			}
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *rootModel) View() string {
	if m.state == loading {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.progress.View())
	}
	return m.list.View()
}

func (m *rootModel) statusCmd(text string) tea.Cmd {
	return func() tea.Msg {
		return statusMsg{text: text}
	}
}

func (m *rootModel) handleUnsubscribeError(mail inbox.MailingList, err error) tea.Cmd {
	var msg string
	if err == inbox.ErrNoUnsubscriber {
		msg = fmt.Sprintf("%s has no one-click unsubscribe feature", mail.From)
	} else {
		msg = fmt.Sprintf("Error unsubscribing from %s: %s", mail.From, err)
	}
	return m.statusCmd(msg)
}
