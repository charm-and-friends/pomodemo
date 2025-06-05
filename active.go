package main

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/v2/progress"
	"github.com/charmbracelet/bubbles/v2/timer"
	tea "github.com/charmbracelet/bubbletea/v2"
	huh "github.com/charmbracelet/huh/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type Session struct {
	rest bool
	// There's no way to access the original timeout duration once the countdown
	// has started.
	duration time.Duration
	timer    timer.Model
	progress progress.Model
	tally    int
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

type (
	ContinueMsg string
	ErrMsg      error
)

func (m Session) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	//	case NewSessionMsg:
	//		m.count = msg.value
	//		return m, nil
	// Reset timer
	// idk when is TimeoutMsg sent/received... it hangs on first load...
	case SettingsMsg:
		// we cannot rest until we've done a full minute of work.
		if m.tally != 0 {
			m.rest = !m.rest
			log.Printf("rest: %v", m.rest)
		}
		if m.rest {
			m.duration, m.err = time.ParseDuration(msg.rest)
		} else {
			m.tally += int(m.duration.Minutes())
			m.duration, m.err = time.ParseDuration(msg.work)
		}
		m.timer = timer.New(m.duration)
		cmds = append(cmds, m.timer.Init())
		return m, tea.Batch(cmds...)
	}
	m.timer, cmd = m.timer.Update(msg)
	return m, cmd
	//	cmds = append(cmds, cmd)
	//	if m.err != nil {
	//		cmds = append(cmds, func() tea.Msg { return ErrMsg(m.err) })
	//	}
	//  return m, tea.Batch(cmds...)
}

func (m Session) View() string {
	message := "Get to work"
	if m.rest {
		message = "Time to snooze..."
	}
	completed := fmt.Sprintf("%d minutes of pure UNBRIDLED focus", m.tally)
	if m.tally == 0 {
		completed = ""
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		completed,
		message,
		// Note: ViewAs is easier than wrangling progress in Update, btw.
		m.progress.ViewAs(float64(m.duration.Seconds()-m.timer.Timeout.Seconds())/m.duration.Seconds()),
		m.timer.View(),
	)
}

/* Confirmation */

// Show next up.
func ContinueMenu() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key("submit").
				Title("Ready?").
				Validate(func(v bool) error {
					if !v {
						return fmt.Errorf("No problem, ready when you are!")
					}
					return nil
				}).
				Negative("Not yet").
				Affirmative("Let's go!"),
		),
	)
}

/* Form */

// TODO
// - work/break timer duration
// - timer title
// - Toggle between textinput vs select

func SettingsMenu() *huh.Form {
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
				Key("save").
				Title("Would you like to save today's sessions?").
				Negative("Don't Save").
				Affirmative("Save"),
			huh.NewConfirm().
				Key("submit").
				Title("Ready?").
				Validate(func(v bool) error {
					if !v {
						return fmt.Errorf("Uh okay then... Figure it out")
					}
					return nil
				}).
				Negative("Nevermind").
				Affirmative("Done"),
		),
	)
}
