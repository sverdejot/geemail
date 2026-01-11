package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	maxWidth = 80

	helpStringTemplate = "You have approximately (%d) mails alongs its submessages. Fetching all of'em to look for mailing list"
)

type progressMsg int64
type endMsg struct{}

type mailLoadingProgress struct {
	progress       progress.Model
	total, current int64
}

func NewProgressModel(total int64) *mailLoadingProgress {
	return &mailLoadingProgress{
		progress: progress.New(progress.WithDefaultGradient()),
		total:    total,
		current:  0,
	}
}

func (m *mailLoadingProgress) Init() tea.Cmd {
	return nil
}

func (m *mailLoadingProgress) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg, endMsg:
		return m, nil

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case mailReceivedMsg:
        m.current++
		cmd := m.progress.SetPercent(float64(m.current) / float64(m.total))
		return m, cmd

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m *mailLoadingProgress) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		titleText+"\n",
		fmt.Sprintf("(%d) ", m.current)+m.progress.View()+"\n",
		fmt.Sprintf(helpStringTemplate, m.total),
	)
}
