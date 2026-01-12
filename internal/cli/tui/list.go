package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sverdejot/geemail/internal/inbox"
)

type viewMode int

const (
	mailingListsOnly viewMode = iota
	allMailGrouped
)

type mailList struct {
	list            list.Model
	viewMode        viewMode
	mailingLists    []inbox.MailingList
	allMailsGrouped []inbox.MailingList
}

func NewModel(mailingLists, allMailsGrouped []inbox.MailingList) mailList {
	items := toListItems(mailingLists)

	listModel := list.New(items, list.NewDefaultDelegate(), 0, 0)
	listModel.Title = "Mailing lists"
	listModel.Styles.Title = titleStyle
	listModel.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			unsubscribe,
			deleteAll,
			archiveAll,
			trashAll,
			markRead,
			markSpam,
			inspect,
			toggleViewMode,
		}
	}
	listModel.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			unsubscribe,
			deleteAll,
			archiveAll,
			trashAll,
			markRead,
			markSpam,
			inspect,
			toggleViewMode,
			toggleHelpMenu,
		}
	}

	return mailList{
		list:            listModel,
		viewMode:        mailingListsOnly,
		mailingLists:    mailingLists,
		allMailsGrouped: allMailsGrouped,
	}
}

func toListItems(mails []inbox.MailingList) []list.Item {
	items := make([]list.Item, 0, len(mails))
	for _, m := range mails {
		items = append(items, m)
	}
	return items
}

func (m mailList) Init() tea.Cmd {
	return nil
}

func (m mailList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil
		case key.Matches(msg, unsubscribe):
			cmds = append(cmds, m.handleUnsubscribe()...)
		case key.Matches(msg, deleteAll):
			cmds = append(cmds, m.handleDeleteAll()...)
		case key.Matches(msg, archiveAll):
			cmds = append(cmds, m.handleArchiveAll()...)
		case key.Matches(msg, trashAll):
			cmds = append(cmds, m.handleTrashAll()...)
		case key.Matches(msg, markRead):
			cmds = append(cmds, m.handleMarkRead()...)
		case key.Matches(msg, markSpam):
			cmds = append(cmds, m.handleMarkSpam()...)
		case key.Matches(msg, inspect):
			cmds = append(cmds, m.handleInspect()...)
		case key.Matches(msg, toggleViewMode):
			m.toggleView()
			return m, nil
		}
	}

	updatedList, cmd := m.list.Update(msg)
	m.list = updatedList
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *mailList) handleUnsubscribe() []tea.Cmd {
	mail, ok := m.getSelectedMail()
	if !ok {
		return nil
	}

	return []tea.Cmd{
		func() tea.Msg {
			return unsubscribeRequestMsg{
				mail: mail,
				idx:  m.list.Index(),
			}
		},
	}
}

func (m *mailList) handleDeleteAll() []tea.Cmd {
	mail, ok := m.getSelectedMail()
	if !ok {
		return nil
	}

	return []tea.Cmd{
		func() tea.Msg {
			return deleteRequestMsg{
				mail: mail,
				idx:  m.list.Index(),
			}
		},
	}
}

func (m *mailList) handleArchiveAll() []tea.Cmd {
	mail, ok := m.getSelectedMail()
	if !ok {
		return nil
	}

	return []tea.Cmd{
		func() tea.Msg {
			return archiveRequestMsg{
				mail: mail,
				idx:  m.list.Index(),
			}
		},
	}
}

func (m *mailList) handleTrashAll() []tea.Cmd {
	mail, ok := m.getSelectedMail()
	if !ok {
		return nil
	}

	return []tea.Cmd{
		func() tea.Msg {
			return trashRequestMsg{
				mail: mail,
				idx:  m.list.Index(),
			}
		},
	}
}

func (m *mailList) handleMarkRead() []tea.Cmd {
	mail, ok := m.getSelectedMail()
	if !ok {
		return nil
	}

	return []tea.Cmd{
		func() tea.Msg {
			return markReadRequestMsg{
				mail: mail,
				idx:  m.list.Index(),
			}
		},
	}
}

func (m *mailList) handleMarkSpam() []tea.Cmd {
	mail, ok := m.getSelectedMail()
	if !ok {
		return nil
	}

	return []tea.Cmd{
		func() tea.Msg {
			return markSpamRequestMsg{
				mail: mail,
				idx:  m.list.Index(),
			}
		},
	}
}

func (m *mailList) handleInspect() []tea.Cmd {
	mail, ok := m.getSelectedMail()
	if !ok {
		return nil
	}

	return []tea.Cmd{
		func() tea.Msg {
			return inspectRequestMsg{mail: mail}
		},
	}
}

func (m *mailList) getSelectedMail() (inbox.MailingList, bool) {
	mail, ok := m.list.SelectedItem().(inbox.MailingList)
	return mail, ok
}

func (m *mailList) toggleView() {
	if m.viewMode == mailingListsOnly {
		m.viewMode = allMailGrouped
		m.list.Title = "All mail (grouped by sender)"
		m.list.SetItems(toListItems(m.allMailsGrouped))
	} else {
		m.viewMode = mailingListsOnly
		m.list.Title = "Mailing lists"
		m.list.SetItems(toListItems(m.mailingLists))
	}
}

func (m mailList) View() string {
	return appStyle.Render(m.list.View())
}
