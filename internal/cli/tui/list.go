package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sverdejot/geemail/internal/core"
)

type mailList struct {
	list list.Model
}

func NewModel(mails []core.MailingList) mailList {
	items := make([]list.Item, len(mails))
	for i, m := range mails {
		items[i] = m
	}

	mailingList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	mailingList.Title = "Mailing lists"
	mailingList.Styles.Title = titleStyle
	mailingList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			unsubscribe,
		}
	}
	mailingList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			unsubscribe,
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
			mail, ok := m.list.SelectedItem().(core.MailingList)
			if !ok {
				break
			}
			cmds = append(cmds, m.list.NewStatusMessage("Unsubcribed from "+mail.From))
		}
	}

	updatedList, cmd := m.list.Update(msg)
	m.list = updatedList
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m mailList) View() string {
	return appStyle.Render(m.list.View())
}
