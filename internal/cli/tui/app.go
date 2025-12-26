package tui

import (
	"context"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sverdejot/geemail/internal/core"
)

type model struct {
	list         list.Model
	svc          *core.MessageService
}

func NewModel(s *core.MessageService) model {
	mails, err := s.GetContent(context.Background())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

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

	return model{
		svc:          s,
		list:         mailingList,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
            mail, ok := m.list.SelectedItem().(core.Content)
            if !ok {
                break
            }
            cmds = append(cmds, m.list.NewStatusMessage("Unsubcribed from " + mail.From))
		}
	}

	updatedList, cmd := m.list.Update(msg)
	m.list = updatedList
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return appStyle.Render(m.list.View())
}
