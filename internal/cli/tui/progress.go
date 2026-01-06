package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	maxWidth = 80
)

type progressMsg int64
type endMsg struct{}

type mailLoadingProgress struct {
	progress progress.Model
    total, current int64
    ch chan struct{}
}

func NewProgressModel(total int64, ch chan struct{}) *mailLoadingProgress {
    return &mailLoadingProgress{
        progress: progress.New(progress.WithDefaultGradient()),
        total: total,
        ch: ch,
    }
}

func (m *mailLoadingProgress) Init() tea.Cmd {
	return m.progressCmd()
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

	case progressMsg:
        m.current += 1
		cmd := m.progress.SetPercent(float64(m.current) / float64(m.total))

		if m.progress.Percent() >= 1.0 {
			return m, m.progressCmd()
		}
		return m, tea.Batch(m.progressCmd(), cmd)

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m *mailLoadingProgress) progressCmd() tea.Cmd {
    return func() tea.Msg {
        _, ok := <-m.ch
        if !ok {
            return endMsg{}
        }
        return progressMsg(1)
    }
}

func (m *mailLoadingProgress) View() string {
	return "\n" +
		m.progress.View() + fmt.Sprintf(" (%d of %d)", m.current, m.total) +"\n\n"
}
