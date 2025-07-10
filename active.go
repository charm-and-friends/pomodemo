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

const (
	work = "work"
	rest = "rest"
)

type Session struct {
	active string
	// There's no way to access the original timeout duration once the countdown
	// has started.
	timer    timer.Model
	progress progress.Model
	tally    int
	err      error

	times map[string]time.Duration
}

func NewSession() Session {
	// create a timer with the desired interval
	times := map[string]time.Duration{
		work: time.Duration(time.Second * 15),
		rest: time.Duration(time.Second * 15),
	}
	session := Session{
		active: work,
		timer:  timer.New(times[work]),
		progress: progress.New(
			progress.WithDefaultGradient(),
			progress.WithWidth(40),
		),
	}
	session.times = times
	return session
}

func (m Session) Init() tea.Cmd {
	return m.timer.Init()
}

type (
	WorkMsg string // Define if we are working or taking a break in the upcoming session.
	MenuMsg struct {
		active string
	}
	ErrMsg error
)

func (m Session) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case timer.TickMsg:
		m.timer, cmd = m.timer.Update(msg)
		cmds = append(cmds, cmd)
	case timer.TimeoutMsg:
		m.Toggle()
		// TODO do we need to pass active
		return m, func() tea.Msg { return MenuMsg{active: m.active} }
	case WorkMsg:
		// Hold onto this for now - might need to pass some info along here at one point.
		return m, m.newTimer()
		//		cmds = append(cmds, m.newTimer())
	case SettingsMsg:
		var err error
		if m.times[rest], err = time.ParseDuration(msg.rest); err != nil {
			m.err = fmt.Errorf("unable to parse rest duration from form %w", err)
		}
		if m.times[work], err = time.ParseDuration(msg.work); err != nil {
			m.err = fmt.Errorf("unable to parse work duration from form %w", err)
		}
		cmds = append(cmds, m.newTimer())
	}
	return m, tea.Sequence(cmds...)
}

func (m *Session) newTimer() tea.Cmd {
	m.timer = timer.New(m.times[m.active])
	return m.timer.Init()
}

func (m Session) View() string {
	message := "Get to work"
	if m.active == rest {
		message = "Time to snooze..."
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		message,
		// Note: ViewAs is easier than wrangling progress in Update, btw.
		m.progress.ViewAs(float64(m.times[m.active].Seconds()-m.timer.Timeout.Seconds())/m.times[m.active].Seconds()),
		m.timer.View(),
	)
}

/* Confirmation */

type ContinueMenu struct {
	active string
	form   *huh.Form
}

func NewContinueMenu(active string) *ContinueMenu {
	return &ContinueMenu{
		active,
		createContinueForm(active),
	}
}

func (c *Session) Toggle() {
	if c.active == rest {
		c.active = work
	} else {
		c.active = rest
	}
}

// Implement tea.Model to be used as a view in our main model.
func (c *ContinueMenu) Init() tea.Cmd {
	return c.form.Init()
}

func (c *ContinueMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c.form.Update(msg)
}

func (c *ContinueMenu) View() string {
	return c.form.View()
}

// Show next up.
func createContinueForm(active string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key("submit").
				Title(fmt.Sprintf("Ready to %s", active)).
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

func SettingsMenu() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
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
