package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/v2/progress"
	"github.com/charmbracelet/bubbles/v2/timer"
	tea "github.com/charmbracelet/bubbletea/v2"
	huh "github.com/charmbracelet/huh/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// TODO make this work on breaks too plz

type Session struct {
	rest bool
	// There's no way to access the original timeout duration once the countdown
	// has started.
	duration time.Duration
	timer    timer.Model
	progress progress.Model
	err      error
}

func NewSession() Session {
	// create a timer with the desired interval
	session := Session{
		timer: timer.New(time.Duration(time.Second * 15)),
		progress: progress.New(
			progress.WithDefaultGradient(),
			progress.WithWidth(40),
		),
	}
	return session
}

func (m Session) Init() tea.Cmd {
	return m.timer.Init()
}

func (m Session) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	// toggle between break and work sessions
	if m.timer.Timedout() {
		m.rest = !m.rest
	}
	switch msg := msg.(type) {
	case SettingsMsg:
		if m.rest {
			m.duration, m.err = time.ParseDuration(msg.rest)
		} else {
			m.duration, m.err = time.ParseDuration(msg.work)
		}
		m.timer = timer.New(m.duration)
		return m, m.timer.Init()
	}
	m.timer, cmd = m.timer.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Session) Done() bool {
	return m.timer.Timedout()
}

func (m Session) View() string {
	message := "Get to work"
	if m.rest {
		message = "Time to snooze..."
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		message,
		// Note: ViewAs is easier than wrangling progress in Update, btw.
		m.progress.ViewAs(float64(m.duration.Seconds()-m.timer.Timeout.Seconds())/m.duration.Seconds()),
		m.timer.View(),
	)
}

/* Form */

// TODO
// - work/break timer duration
// - timer title
// - Toggle between textinput vs select

func NewForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			// enter work time
			// enter break time
			// confirm
			// show same form if I'm modifying, but set placeholder text to current values
			huh.NewSelect[string]().
				Key("work").
				Options(huh.NewOptions("10s", "25m", "30m", "45m", "50m", "60m")...).
				Title("Set work session duration").
				Description("Get to work ya little"),
			huh.NewSelect[string]().
				Key("rest").
				Options(huh.NewOptions("5s", "2m", "5m", "10m", "15m")...).
				Title("Set rest session duration").
				Description("Nap time"),
			huh.NewConfirm().
				Key("submit").
				Title("Ready?").
				Validate(func(v bool) error {
					if !v {
						return fmt.Errorf("Uh okay then... Figure it out")
					}
					return nil
				}).
				Affirmative("Done").
				Negative("Nevermind"),
		),
	)
}
