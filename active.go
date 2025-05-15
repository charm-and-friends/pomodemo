package main

import (
	"time"

	"github.com/charmbracelet/bubbles/v2/progress"
	"github.com/charmbracelet/bubbles/v2/timer"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type active struct {
	timer    timer.Model
	progress progress.Model
}

func NewActive() active {
	// create a timer with the desired interval
	return active{
		timer:    timer.New(time.Duration(time.Second * 15)),
		progress: progress.New(),
	}
}

func (m active) Init() tea.Cmd {
	return m.timer.Init()
}

func (m active) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	// switch msg := msg.(type) {
	// case timer.TickMsg:

	// }
	m.timer, cmd = m.timer.Update(msg)
	return m, cmd
}

func (m active) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		"Get to work",
		m.timer.View(),
	)
}
