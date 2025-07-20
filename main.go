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
	isReady   bool // 是否输入任务名
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
		// 每秒都会触发，不管是暂停还是运行
		if m.isRunning {
			// 正常倒计时
			if m.timeLeft > 0 {
				m.timeLeft -= time.Second
				if m.timeLeft <= 0 {
					return m, tea.Quit
				}
			}
		}
		// 如果暂停或时间到了，不再返回 tick
		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "s":
			// 只切换状态，不注册新 tick（避免多重加速）
			if m.isReady {
				m.isRunning = !m.isRunning
				return m, nil
			}
		}

		if !m.isReady {
			// 阶段 1：任务名输入
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
	// 阶段 1：任务名输入
	if !m.isReady {
		style := lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			Foreground(lipgloss.Color("12")).
			BorderForeground(lipgloss.Color("228")).
			Width(40)

		return style.Render("Enter task name: \n\n" + renderInput(m.textInput))
	}

	// 倒计时逻辑
	total := int(m.timeLeft.Seconds())
	if total < 0 {
		total = 0
	}
	mins := total / 60
	secs := total % 60

	// 基础内容
	content := fmt.Sprintf(
		"🍅 Task: %s\n\n⏳ Time Left：%02dm:%02ds",
		m.taskName, mins, secs,
	)

	// 如果暂停，则插入“暂停提示”
	if !m.isRunning {
		pauseBox := lipgloss.NewStyle().
			Foreground(lipgloss.Color("13")).
			Background(lipgloss.Color("235")).
			Padding(0, 1).
			Bold(true).
			Render("⏸ Paused, Enter [s] to continue")

		// 插入到 content中（加在倒计时下面）
		content += "\n\n" + pauseBox
	}

	// 控制提示
	controls := "\n\n[s] Start/Pause	[q] Quit"

	// 整体框样式
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
	// 初始化textInput并写入model
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
