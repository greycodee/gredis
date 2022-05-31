package views

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/greycodee/gredis/core"
	"strings"
)

type mainPage struct {
	cmdInput textinput.Model
	history	[]string
	redisCli *core.RedisClient
	cmdResult string
}

func initMainPage(redisCli *core.RedisClient) mainPage {
	t := textinput.New()
	t.Focus()
	h := []string{"select 0","get name"}
	mp := mainPage{cmdInput: t,history: h,redisCli: redisCli}
	return mp
}

func (mp mainPage) Init()  tea.Cmd{
	mp.cmdInput = textinput.New()
	mp.cmdInput.Focus()
	return nil
}

func (mp mainPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return mp, tea.Quit
		case "enter":
			if len(mp.history)+1>10 {
				mp.history = mp.history[1:]
			}
			mp.cmdResult = string(mp.redisCli.ExecCMD(mp.cmdInput.Value()))
			mp.history = append(mp.history, mp.cmdInput.Value())
			mp.cmdInput.Reset()

			// TODO 执行命令
			return mp,nil

		}
	}
	cmd := mp.updateInputs(msg)
	return mp,cmd
}

func (mp mainPage) View() string {
	return mp.mainView()
}

func (mp *mainPage) updateInputs(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	mp.cmdInput,cmd = mp.cmdInput.Update(msg)

	return tea.Batch(cmd)
}

func (mp *mainPage) mainView() string {
	status := lipgloss.NewStyle().
		MarginLeft(1).
		MarginRight(5).
		Padding(0, 1).
		Italic(true).
		Height(3).
		Width(50).
		Border(lipgloss.NormalBorder()).
		Foreground(lipgloss.Color("#eedd22")).
		SetString("Address:localhost.com:6379\nUsername:root\nDatabases:0")

	h := strings.Builder{}
	for _,v := range mp.history{
		h.WriteString(v)
		h.WriteString("\n")
	}
	cmdHistory := lipgloss.NewStyle().
		MarginLeft(1).
		MarginRight(5).
		Padding(0, 1).
		Italic(true).
		Width(50).
		Border(lipgloss.NormalBorder()).
		Foreground(lipgloss.Color("#eedd22")).
		SetString(h.String())

	cmd := lipgloss.NewStyle().
		MarginLeft(1).
		MarginRight(5).
		Padding(0, 1).
		Italic(true).
		Height(1).
		Width(50).
		Border(lipgloss.NormalBorder()).
		Foreground(lipgloss.Color("#eedd22")).
		SetString(mp.cmdInput.View())

	content := lipgloss.NewStyle().
		MarginLeft(1).
		MarginRight(5).
		Padding(0, 1).
		Italic(true).
		Height(10).
		Width(50).
		Border(lipgloss.NormalBorder()).
		Foreground(lipgloss.Color("#eedd22")).
		SetString(mp.cmdResult)


	mm := lipgloss.JoinVertical(lipgloss.Left,status.String(),cmdHistory.String(),cmd.String())
	m2 := lipgloss.JoinHorizontal(lipgloss.Left,mm,content.String())
	return m2
}
