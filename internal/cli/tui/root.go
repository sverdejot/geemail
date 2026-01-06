package tui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sverdejot/geemail/internal/core"
)

type state int

const (
	loading state = iota
	ready
)

type rootModel struct {
	state    state
	progress *mailLoadingProgress
	list     mailList
	ctx      context.Context
	svc      *core.MailService
	mails    []core.RawMail
	inc      chan struct{}
	width    int
	height   int
}

func NewRoot(ctx context.Context, svc *core.MailService) (*rootModel, error) {
	total, err := svc.GetTotalUnreads(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get total unread messages: %w", err)
	}

	inc := make(chan struct{})
	mails := make([]core.RawMail, 0)

	pg := NewProgressModel(total, inc)

	return &rootModel{
		state:    loading,
		progress: pg,
		ctx:      ctx,
		svc:      svc,
		mails:    mails,
		inc:      inc,
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
			return tea.Quit()
		}

		go func() {
			for mail := range stream {
				m.mails = append(m.mails, mail)
				select {
				case m.inc <- struct{}{}:
				case <-m.ctx.Done():
					return
				}
			}
			close(m.inc)
		}()

		return nil
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

	case progressMsg:
		if m.state == loading {
			_, cmd := m.progress.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

	case endMsg:
		if m.state == loading {
			rawMailList := core.RawMailList(m.mails)
			mailingLists := core.GetMailingList(rawMailList)
			
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
		return m.progress.View()
	}
	return m.list.View()
}

