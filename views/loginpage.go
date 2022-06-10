package views

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/greycodee/gredis/core"
	"os"
	"strconv"
	"strings"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))

	defaultPort			= "6379"
	defaultDatabases	= "0"
)

type loginPage struct {
	focusIndex int
	inputs     []textinput.Model
	serverAddr  string
	password    string
	serverPort  string
	databases   string
}

func InitialModel() loginPage {
	m := loginPage{
		inputs: make([]textinput.Model, 4),
		serverAddr: "127.0.0.1",
		serverPort: "6379",
		databases: "0",
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Host"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Port      [default:6379]"

		case 2:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		case 3:
			t.Placeholder = "Databases [default:0]"
			t.CharLimit=2
		}

		m.inputs[i] = t
	}

	return m
}
func (m loginPage) Init() tea.Cmd {
	return textinput.Blink
}

func (m loginPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				// 进行提交判断
				if m.inputs[0].Value() == ""{
					m.focusIndex = 0
					m.inputs[0].Placeholder="Please enter redis server host!"
					return m.focusesUpdate(nil)
				}else if i,e := strconv.ParseInt(m.inputs[3].Value(),10,32);e != nil || i > 60 || i<0{
					//
					m.focusIndex = 3
					m.inputs[3].SetValue("")
					m.inputs[3].Placeholder="Please enter 0-60 databases!"
					return m.focusesUpdate(nil)
				}else {
					redisClient := &core.RedisClient{}
					err := redisClient.Open(m.serverAddr + ":" + m.serverPort)
					if err != nil {
						m.focusIndex = 0
						m.inputs[0].SetValue("")
						m.inputs[0].Placeholder="Host or Port error!"
						m.inputs[1].SetValue("")
						m.inputs[1].Placeholder="Host or Port error!"
						return m.focusesUpdate(nil)
					}
					// 登陆认证 auth
					if m.password != ""{
						loginResp := redisClient.ExecCMD("auth",m.password)
						if string(loginResp)!="OK" {
							// 登陆失败
							m.focusIndex = 2
							m.inputs[2].SetValue("")
							m.inputs[2].Placeholder="Password error!"
							return m.focusesUpdate(nil)
						}
					}else{
						p2 := tea.NewProgram(initMainPage(redisClient,m.serverAddr,m.serverPort,m.databases),tea.WithAltScreen())
						if err := p2.Start(); err != nil {
							fmt.Printf("could not start program: %s\n", err)
							os.Exit(1)
						}
						return m, tea.Quit
					}

				}

			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			return m.focusesUpdate(nil)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m loginPage) focusesUpdate(cmds []tea.Cmd) (tea.Model, tea.Cmd) {
	if m.focusIndex > len(m.inputs) {
		m.focusIndex = 0
	} else if m.focusIndex < 0 {
		m.focusIndex = len(m.inputs)
	}
	innerCmds := make([]tea.Cmd, len(m.inputs))
	for i := 0; i <= len(m.inputs)-1; i++ {
		if i == m.focusIndex {
			// Set focused state
			innerCmds[i] = m.inputs[i].Focus()
			m.inputs[i].PromptStyle = focusedStyle
			m.inputs[i].TextStyle = focusedStyle
			continue
		}
		// Remove focused state
		m.inputs[i].Blur()
		m.inputs[i].PromptStyle = noStyle
		m.inputs[i].TextStyle = noStyle
	}
	innerCmds = append(innerCmds,cmds...)
	return m, tea.Batch(innerCmds...)
}

func (m *loginPage) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	m.serverAddr = m.inputs[0].Value()
	if m.inputs[1].Value() != "" {
		m.serverPort = m.inputs[1].Value()
	}else{
		m.serverPort = defaultPort
	}
	m.password = m.inputs[2].Value()
	if m.inputs[3].Value() != "" {
		m.databases = m.inputs[3].Value()
	}else {
		m.databases = defaultDatabases
	}
	return tea.Batch(cmds...)
}

func (m loginPage) View() string {

		var b strings.Builder
		for i := range m.inputs {
			b.WriteString(m.inputs[i].View())
			if i < len(m.inputs)-1 {
				b.WriteRune('\n')
			}
		}

		button := &blurredButton
		if m.focusIndex == len(m.inputs) {
			button = &focusedButton
		}
		fmt.Fprintf(&b, "\n\n%s\n\n", *button)

		b.WriteString(helpStyle.Render("hi, welcome come to use gredis!"))

		return b.String()

}
