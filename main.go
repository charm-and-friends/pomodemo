package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea/v2"
	huh "github.com/charmbracelet/huh/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/dgraph-io/badger"
)

// TODO create pomodoro app multiple views
// - form to define pomodoro pref (work session, break, maybe even a goal?)
// - active session (timer countdown with progress bar, show how many sessions
// they've done, maybe even textarea or list to track tasks?)
//   - active session commands:
//   - pause, stop, skip, quit, restart

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal("fatal:", err)
	}
	defer f.Close()

	p := tea.NewProgram(NewMain())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

const (
	settings int = iota
	session
	confirm
)

/* Main Model  */

type Model struct {
	active   int
	views    []tea.Model
	form     *huh.Form
	settings SettingsMsg
	db       *badger.DB
	err      error
}

func NewMain() *Model {
	return &Model{
		views: []tea.Model{SettingsMenu(), NewSession(), NewContinueMenu(work)},
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

type NewSessionMsg struct {
	err   error
	value string
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.views[m.active], cmd = m.views[m.active].Update(msg)
	cmds = append(cmds, cmd)
	if form, ok := m.views[m.active].(*huh.Form); ok {
		if form.State == huh.StateCompleted {
			// set values
			if m.active == settings {
				m.settings = SettingsMsg{
					work: form.GetString("work"),
					rest: form.GetString("rest"),
				}
				m.views[session], cmd = m.views[session].Update(m.settings)
				m.active = session
				return m, cmd
			}
			if m.active == confirm {
				// Whether or not we just completed a working session.
				m.active = session
				m.views[m.active], cmd = m.views[m.active].Update(WorkMsg(""))
				return m, cmd
			}
		}
	}

	switch msg := msg.(type) {
	// Confirm we're ready for next session before moving forward.
	case MenuMsg:
		m.views[confirm] = NewContinueMenu(msg.active)
		m.active = confirm
		return m, m.views[m.active].Init()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case ErrMsg:
		m.err = msg
	}

	return m, tea.Sequence(cmds...)
}

func (m Model) View() string {
	var err string
	if m.err != nil {
		err = m.err.Error()
	}
	if view, ok := m.views[m.active].(tea.ViewModel); ok {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			view.View(),
			err,
		)
	}
	return "no view models :("
}
