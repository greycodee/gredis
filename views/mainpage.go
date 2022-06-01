package views

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/greycodee/gredis/core"
	"github.com/greycodee/gredis/core/cmd"
	"strings"
)

type mainPage struct {
	cmdInput textinput.Model
	history	[]string
	redisCli *core.RedisClient
	cmdResult string
	serverAddr string
	serverPort string
	databases string
}

func initMainPage(redisCli *core.RedisClient,addr string,port string, databases string,) mainPage {
	t := textinput.New()
	t.Focus()
	h := make([]string,0)
	mp := mainPage{
		cmdInput: t,
		history: h,
		redisCli: redisCli,
		serverAddr: addr,
	serverPort: port,
	databases: databases}
	return mp
}

func (mp mainPage) Init()  tea.Cmd{
	mp.cmdInput = textinput.New()
	mp.cmdInput.Focus()
	mp.redisCli.ExecCMD("select",mp.databases)
	return nil
}

func (mp mainPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			mp.redisCli.Shutdown()
			return mp, tea.Quit
		case "enter":
			if len(mp.history)+1>10 {
				mp.history = mp.history[1:]
			}
			cmdByte,selectDB,db := cmd.GetCmdByte(mp.cmdInput.Value())

			mp.cmdResult = string(mp.redisCli.ExecCMDByte(cmdByte))
			if selectDB && mp.cmdResult=="OK"{
				mp.databases = db
			}
			mp.history = append(mp.history, mp.cmdInput.Value())
			mp.cmdInput.Reset()

			return mp,nil

		}
	}
	inputs := mp.updateInputs(msg)
	return mp, inputs
}

func (mp mainPage) View() string {
	return mp.mainView()
}

func (mp *mainPage) updateInputs(msg tea.Msg) tea.Cmd {
	var c tea.Cmd
	mp.cmdInput, c = mp.cmdInput.Update(msg)

	return tea.Batch(c)
}

func (mp *mainPage) mainView() string {
	statusContent := fmt.Sprintf("Address: %s\nPort: %s\nDatabases: %s",mp.serverAddr,mp.serverPort,mp.databases)
	status := lipgloss.NewStyle().
		MarginLeft(1).
		MarginRight(5).
		Padding(0, 1).
		Italic(true).
		Height(3).
		Width(50).
		Border(lipgloss.NormalBorder()).
		Foreground(lipgloss.Color("#eedd22")).
		SetString(statusContent)

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

	cmdInput := lipgloss.NewStyle().
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


	mm := lipgloss.JoinVertical(lipgloss.Left,status.String(),cmdHistory.String(),cmdInput.String())
	m2 := lipgloss.JoinHorizontal(lipgloss.Left,mm,content.String())
	return m2
}
