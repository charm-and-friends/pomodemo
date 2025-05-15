package main

import (
	"time"

	"github.com/charmbracelet/bubbles/v2/progress"
	"github.com/charmbracelet/bubbles/v2/timer"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// TODO make this work on breaks too plz

type active struct {
	timer    timer.Model
	progress progress.Model
}

func NewActive() active {
	// create a timer with the desired interval
	return active{
		timer: timer.New(time.Duration(time.Second * 15)),
		progress: progress.New(
			progress.WithDefaultGradient(),
			progress.WithWidth(40),
		),
	}
}

func (m active) Init() tea.Cmd {
	return m.timer.Init()
}

func (m active) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.timer, cmd = m.timer.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m active) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		"Get to work",
		// Note: ViewAs is easier than wrangling progress in Update, btw.
		m.progress.ViewAs(float64(15-m.timer.Timeout.Seconds())/15),
		m.timer.View(),
	)
}
