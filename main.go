package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea/v2"
	huh "github.com/charmbracelet/huh/v2"
)

// TODO create pomodoro app multiple views
// - form to define pomodoro pref (work session, break, maybe even a goal?)
// - active session (timer countdown with progress bar, show how many sessions
// they've done, maybe even textarea or list to track tasks?)
// 	- active session commands:
// 	 	- pause, stop, skip, quit, restart

func main() {
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

const (
	form int = iota
	session
)

type Model struct {
	// active, new, history? break?
	active   int
	views    []tea.Model
	form     *huh.Form
	settings SettingsMsg
}

func NewModel() *Model {
	return &Model{
		views: []tea.Model{NewForm(), NewSession()},
	}
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, view := range m.views {
		cmds = append(cmds, view.Init())
	}
	return tea.Batch(cmds...)
}

type SettingsMsg struct {
	work string
	rest string
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.views[m.active], cmd = m.views[m.active].Update(msg)
	// If we're in the form view, let's check when the user has submitted.
	if form, ok := m.views[m.active].(*huh.Form); ok {
		if form.State == huh.StateCompleted {
			m.settings = SettingsMsg{
				work: form.GetString("work"),
				rest: form.GetString("rest"),
			}
			m.views[session], cmd = m.views[session].Update(m.settings)
			m.active = session
			return m, cmd
		}
	}

	if session, ok := m.views[m.active].(Session); ok {
		if session.Done() {
			m.views[m.active], cmd = m.views[m.active].Update(m.settings)
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m Model) View() string {
	if view, ok := m.views[m.active].(tea.ViewModel); ok {
		return view.View()
	}
	return "no view models :("
}
