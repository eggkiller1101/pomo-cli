package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	textInput textinput.Model
	timeLeft  time.Duration
	isRunning bool
	isReady   bool // if the user has entered a task name
	taskName  string
}

type tickMsg time.Time

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tickMsg:
		// Every second will trigger, whether paused or running
		if m.isRunning {
			// Normal countdown
			if m.timeLeft > 0 {
				m.timeLeft -= time.Second
				if m.timeLeft <= 0 {
					return m, tea.Quit
				}
			}
		}
		// If paused or time's up, don't return tick
		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "s":
			// Only switch state, don't register new tick (avoid multiple acceleration)
			if m.isReady {
				m.isRunning = !m.isRunning
				return m, nil
			}
		}

		if !m.isReady {
			// Stage 1: Task name input
			switch msg.Type {
			case tea.KeyEnter:
				m.taskName = m.textInput.Value()
				if m.taskName != "" {
					m.isReady = true
					m.isRunning = true
					return m, tick() // 首次启动
				}
				return m, nil
			default:
				var cmd tea.Cmd
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd
			}
		}
	}

	return m, nil
}

func renderInput(ti textinput.Model) string {
	if ti.Value() == "" {
		return "> _"
	}
	return ti.View()
}

func (m model) View() string {
	// Stage 1: Task name input
	if !m.isReady {
		style := lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			Foreground(lipgloss.Color("12")).
			BorderForeground(lipgloss.Color("228")).
			Width(40)

		return style.Render("Enter task name: \n\n" + renderInput(m.textInput))
	}

	// Countdown logic
	total := int(m.timeLeft.Seconds())
	if total < 0 {
		total = 0
	}
	mins := total / 60
	secs := total % 60

	// Basic content
	content := fmt.Sprintf(
		"🍅 Task: %s\n\n⏳ Time Left：%02dm:%02ds",
		m.taskName, mins, secs,
	)

	// If paused, insert "paused" prompt
	if !m.isRunning {
		pauseBox := lipgloss.NewStyle().
			Foreground(lipgloss.Color("13")).
			Background(lipgloss.Color("235")).
			Padding(0, 1).
			Bold(true).
			Render("⏸ Paused, Enter [s] to continue")

		// Insert into content (below the countdown)
		content += "\n\n" + pauseBox
	}

	// Controls
	controls := "\n\n[s] Start/Pause	[q] Quit"

	// Overll layout style
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("228")).
		Padding(1, 2).
		Align(lipgloss.Center).
		Width(40).
		Foreground(lipgloss.Color("12"))

	return style.Render(content + controls)
}

func main() {
	// Initialize textInput and write it in model
	ti := textinput.New()
	ti.Placeholder = "Enter your task name"
	ti.Focus()
	ti.Cursor.Style = lipgloss.NewStyle().
		Underline(true).
		Foreground(lipgloss.Color("12"))

	m := model{
		textInput: ti,
		timeLeft:  10 * time.Minute,
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
