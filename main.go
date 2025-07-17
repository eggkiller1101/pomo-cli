package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	timeLeft time.Duration
}

type tickMsg time.Time

func (m model) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tickMsg:
		if m.timeLeft <= 0 {
			return m, tea.Quit
		}
		m.timeLeft -= time.Second
		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	total := int(m.timeLeft.Seconds())
	if total < 0 {
		total = 0
	}
	mins := total / 60
	secs := total % 60
	return fmt.Sprintf("⏳ Time left: %02dm:%02ds\nPress q to quit.\n", mins, secs)
}

func main() {
	m := model{timeLeft: 10 * time.Minute}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
