package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// TODO create pomodoro app multiple views
// - form to define pomodoro pref (work session, break, maybe even a goal?)
// - active session (timer countdown with progress bar, show how many sessions
// they've done, maybe even textarea or list to track tasks?)
// 	- active session commands:
// 	 	- pause, stop, skip, quit, restart

func main() {
	fmt.Println("hello world")
}

type Model struct{}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return "hello model"
}
