package tui

import (
	"fmt"
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sverdejot/geemail/internal/inbox"
)

// mailItem wraps RawMail to display Subject as title instead of From
type mailItem struct {
	mail inbox.RawMail
}

func (m mailItem) FilterValue() string {
	return m.mail.Subject
}

func (m mailItem) Title() string {
	return m.mail.Subject
}

func (m mailItem) Description() string {
	return m.mail.Snippet
}

type inspectModel struct {
	list   list.Model
	sender string
}

func NewInspectModel(sender string, mails []inbox.RawMail) inspectModel {
	// Sort by date descending (parse Date header)
	sortByDateDesc(mails)

	items := make([]list.Item, 0, len(mails))
	for _, m := range mails {
		items = append(items, mailItem{mail: m})
	}

	mailList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	mailList.Title = fmt.Sprintf("Mails from %s", sender)
	mailList.Styles.Title = titleStyle
	mailList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			deleteAll,
			archiveAll,
			trashAll,
			markRead,
			markSpam,
			goBack,
		}
	}
	mailList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			deleteAll,
			archiveAll,
			trashAll,
			markRead,
			markSpam,
			goBack,
			toggleHelpMenu,
		}
	}

	return inspectModel{
		list:   mailList,
		sender: sender,
	}
}

func (m inspectModel) Init() tea.Cmd {
	return nil
}

func (m inspectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case key.Matches(msg, goBack):
			return m, func() tea.Msg { return goBackMsg{} }
		case key.Matches(msg, toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil
		case key.Matches(msg, deleteAll):
			cmds = append(cmds, m.handleDelete()...)
		case key.Matches(msg, archiveAll):
			cmds = append(cmds, m.handleArchive()...)
		case key.Matches(msg, trashAll):
			cmds = append(cmds, m.handleTrash()...)
		case key.Matches(msg, markRead):
			cmds = append(cmds, m.handleMarkRead()...)
		case key.Matches(msg, markSpam):
			cmds = append(cmds, m.handleMarkSpam()...)
		}
	}

	updatedList, cmd := m.list.Update(msg)
	m.list = updatedList
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m inspectModel) View() string {
	return appStyle.Render(m.list.View())
}

func (m *inspectModel) getSelectedMail() (inbox.RawMail, bool) {
	item, ok := m.list.SelectedItem().(mailItem)
	if !ok {
		return inbox.RawMail{}, false
	}
	return item.mail, true
}

func (m *inspectModel) handleDelete() []tea.Cmd {
	mail, ok := m.getSelectedMail()
	if !ok {
		return nil
	}

	return []tea.Cmd{
		func() tea.Msg {
			return singleDeleteRequestMsg{
				mail: mail,
				idx:  m.list.Index(),
			}
		},
	}
}

func (m *inspectModel) handleArchive() []tea.Cmd {
	mail, ok := m.getSelectedMail()
	if !ok {
		return nil
	}

	return []tea.Cmd{
		func() tea.Msg {
			return singleArchiveRequestMsg{
				mail: mail,
				idx:  m.list.Index(),
			}
		},
	}
}

func (m *inspectModel) handleTrash() []tea.Cmd {
	mail, ok := m.getSelectedMail()
	if !ok {
		return nil
	}

	return []tea.Cmd{
		func() tea.Msg {
			return singleTrashRequestMsg{
				mail: mail,
				idx:  m.list.Index(),
			}
		},
	}
}

func (m *inspectModel) handleMarkRead() []tea.Cmd {
	mail, ok := m.getSelectedMail()
	if !ok {
		return nil
	}

	return []tea.Cmd{
		func() tea.Msg {
			return singleMarkReadRequestMsg{
				mail: mail,
				idx:  m.list.Index(),
			}
		},
	}
}

func (m *inspectModel) handleMarkSpam() []tea.Cmd {
	mail, ok := m.getSelectedMail()
	if !ok {
		return nil
	}

	return []tea.Cmd{
		func() tea.Msg {
			return singleMarkSpamRequestMsg{
				mail: mail,
				idx:  m.list.Index(),
			}
		},
	}
}

func sortByDateDesc(mails []inbox.RawMail) {
	sort.Slice(mails, func(i, j int) bool {
		dateI := parseDate(mails[i].Headers["date"])
		dateJ := parseDate(mails[j].Headers["date"])
		return dateI.After(dateJ)
	})
}

func parseDate(dates []string) time.Time {
	if len(dates) == 0 {
		return time.Time{}
	}

	// Common email date formats
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"Mon, 2 Jan 2006 15:04:05 MST",
		"2 Jan 2006 15:04:05 -0700",
		"Mon, 02 Jan 2006 15:04:05 -0700",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dates[0]); err == nil {
			return t
		}
	}

	return time.Time{}
}
