package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sverdejot/geemail/internal/inbox"
)

type mailList struct {
	list list.Model
}

func NewModel(mails []inbox.MailingList) mailList {
	items := make([]list.Item, 0, len(mails))
	for _, m := range mails {
		if !m.UnsubscribeAvailable() {
			continue
		}
		items = append(items, m)
	}

	mailingList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	mailingList.Title = "Mailing lists"
	mailingList.Styles.Title = titleStyle
	mailingList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			unsubscribe,
			deleteAll,
		}
	}
	mailingList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			unsubscribe,
			deleteAll,
			toggleHelpMenu,
		}
	}

	return mailList{
		list: mailingList,
	}
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

func (m *mailList) getSelectedMail() (inbox.MailingList, bool) {
	mail, ok := m.list.SelectedItem().(inbox.MailingList)
	return mail, ok
}

func (m mailList) View() string {
	return appStyle.Render(m.list.View())
}
