package main

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/v2/timer"
	tea "github.com/charmbracelet/bubbletea/v2"
	huh "github.com/charmbracelet/huh/v2"
	"github.com/charmbracelet/log"
	"github.com/dgraph-io/badger"
)

// TODO create pomodoro app multiple views
// - form to define pomodoro pref (work session, break, maybe even a goal?)
// - active session (timer countdown with progress bar, show how many sessions
// they've done, maybe even textarea or list to track tasks?)
// 	- active session commands:
// 	 	- pause, stop, skip, quit, restart
//
//
// TODO confirmation before moving onto break/work.
// would like for this confirmation dialogue to show current number of sessions for the day

func main() {
	//	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	defer db.Close()
	//
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Error("fatal:", err)
	}
	defer f.Close()

	p := tea.NewProgram(NewModel())
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
	// active, new, history? break?
	active   int
	views    []tea.Model
	form     *huh.Form
	settings SettingsMsg
	db       *badger.DB
}

// func NewModel(db *badger.DB) *Model {
// 	return &Model{
// 		db:    db,
// 		views: []tea.Model{SettingsMenu(), NewSession(), ContinueMenu()},
// 	}
// }

func NewModel() *Model {
	return &Model{
		views: []tea.Model{SettingsMenu(), NewSession(), ContinueMenu()},
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

func (m Model) showCount() string {
	var count string
	m.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(m.settings.work))
		count = item.String()
		if err != nil {
			return fmt.Errorf("unable to get number of work sessions")
		}
		return txn.Commit()
	})
	return count
}

// key should be today's date.
func (m Model) increment() tea.Msg {
	msg := NewSessionMsg{
		value: m.showCount(),
	}
	var count int
	count, msg.err = strconv.Atoi(msg.value)

	return m.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(m.settings.work), []byte(fmt.Sprintf("%d", count+1)))
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.views[m.active], cmd = m.views[m.active].Update(msg)
	// If we're in the form view, let's check when the user has submitted.
	if form, ok := m.views[m.active].(*huh.Form); ok {
		if form.State == huh.StateCompleted {
			// set values
			if m.active == settings {
				m.settings = SettingsMsg{
					work: form.GetString("work"),
					rest: form.GetString("rest"),
				}
			}

			m.views[session], cmd = m.views[session].Update(m.settings)
			m.active = session
			return m, cmd

		}
	}

	//	if session, ok := m.views[m.active].(Session); ok {
	//		if session.Done() {
	//			m.views[m.active], cmd = m.views[m.active].Update(m.settings)
	//		}
	//	}

	switch msg := msg.(type) {
	// Confirm we're ready for next session before moving forward.
	case timer.TimeoutMsg:
		m.active = confirm
		log.Print("continue msg received in main")
		m.views[m.active] = ContinueMenu()
		m.views[m.active], cmd = m.views[m.active].Update(msg)
		//	case DoneSessionMsg:
		//		return m, m.increment
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m Model) View() string {
	// TODO show errors
	if view, ok := m.views[m.active].(tea.ViewModel); ok {
		return view.View()
	}
	return "no view models :("
}
