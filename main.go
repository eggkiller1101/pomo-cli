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

type tickMsg struct{}

func (m model) Init() tea.Cmd {
	return tea.Tick(time.Minute, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if m.timeLeft <= 0 {
			return m, tea.Quit
		}
		m.timeLeft -= time.Second
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg{}
		})
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	totalTime := int(m.timeLeft.Seconds())
	mins := totalTime / 60
	secs := totalTime % 60
	return fmt.Sprintf("Time left: %02dm:%02ds\nPress q to quit.", mins, secs)
}

func main() {
	m := model{timeLeft: 10 * time.Minute}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting program:", err)
		os.Exit(1)
	}
}
