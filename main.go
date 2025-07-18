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
		if m.timeLeft <= 0 {
			return m, tea.Quit
		}
		m.timeLeft -= time.Second
		return m, tick()

	case tea.KeyMsg:
		// 退出快捷键设置成全局，任何情况下都能用
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if !m.isReady {
			// 阶段 1：任务名输入
			switch msg.Type {
			case tea.KeyEnter:
				m.taskName = m.textInput.Value()
				if m.taskName != "" {
					m.isReady = true
					m.isRunning = true
					return m, tick()
				}
				return m, nil // 不允许空任务名
			default:
				var cmd tea.Cmd
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd
			}
		} else {
			// 阶段 2：番茄倒计时界面
			switch msg.String() {
			case "s":
				m.isRunning = !m.isRunning
				if m.isRunning {
					return m, tick()
				}
				return m, nil
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	// 阶段 1：任务名输入
	if !m.isReady {
		style := lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			Foreground(lipgloss.Color("12")).
			BorderForeground(lipgloss.Color("228")).
			Width(40)

		return style.Render("请输入任务名后按 Enter 开始：\n\n" + m.textInput.View())
	}
	// 倒计时逻辑
	total := int(m.timeLeft.Seconds())
	if total < 0 {
		total = 0
	}
	mins := total / 60
	secs := total % 60

	// 倒计时界面
	content := fmt.Sprintf(
		"🍅 Task: %s\n\n⏳ Time Left：%02dm:%02ds\n\n[s] Start/Pause   [q] Quit",
		m.taskName, mins, secs,
	)

	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("228")).
		Padding(1, 2).
		Align(lipgloss.Center).
		Width(40).
		Foreground(lipgloss.Color("12"))

	return style.Render(content)
}

func main() {
	// 初始化textInput并写入model
	ti := textinput.New()
	ti.Placeholder = "Enter your task name"
	ti.Focus()

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
