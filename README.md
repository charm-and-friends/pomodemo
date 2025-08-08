# Lapin

A terminal-based Pomodoro timer application built with Go and Charm's TUI libraries.

## Features

- **Customizable Sessions**: Set work and break durations (10s-60m for work, 5s-15m for breaks)
- **Visual Progress**: Progress bar and countdown timer display
- **Session Management**: Automatic transitions between work and break sessions
- **Interactive Forms**: Easy-to-use terminal interface for configuration
- **Session Confirmation**: Prompts before starting each new session

## Installation

```bash
go build -o lapin
./lapin
```

## Usage

1. **Configure Settings**: Set your preferred work and break session durations
2. **Start Working**: Begin your first work session with the countdown timer
3. **Take Breaks**: Automatically prompted for break sessions after work periods
4. **Continue**: Confirm when you're ready to start each new session

## Controls

- `q` or `Ctrl+C`: Quit the application
- Use arrow keys and Enter to navigate forms
- Follow on-screen prompts for session management

## Built With

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Huh](https://github.com/charmbracelet/huh) - Interactive forms
- [Bubbles](https://github.com/charmbracelet/bubbles) - Timer and progress components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling and layout
- [Badger](https://github.com/dgraph-io/badger) - Embedded database (planned for session persistence)

## Development Status

This is a work-in-progress Pomodoro application. Planned features include:
- Session persistence and statistics
- Task tracking during sessions
- Enhanced session controls (pause, skip, restart)
- Session counter and progress tracking