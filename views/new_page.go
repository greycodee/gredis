package views

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/greycodee/gredis/core"
	"github.com/greycodee/gredis/core/cmd"
	"github.com/rivo/tview"
	"strings"
)

type TUI struct {
	serverInfoView  *tview.TextView
	cmdHistoryView	*tview.TextView
	cmdInputView	*tview.InputField
	keyInfoView		*tview.TextView
	keyValueView	*tview.TextView
	keyHexValueView	*tview.TextView

	mainPage 		*tview.Flex

	RedisServer 	RedisServer
	KeyInfo			KeyInfo
	TUIData			TUIData

	focusIndex		int
	allWidgets		[]tview.Primitive
	historyIndex	int
}

type RedisServer struct {
	Addr 	string
	Port 	string
	DB 		string
	Conn 	core.RedisClient
}

type KeyInfo struct {
	TTL 	int64
	Key 	string
	KeyType KeyType
}

type KeyType uint8

const  (
	STRING 	KeyType = 1
	LIST	KeyType = 2
	SET		KeyType = 3
	ZSET	KeyType = 4
	HASH	KeyType = 5
)

type TUIData struct {
	CmdHistory	[]string
	KeyValue 	string
	KeyHexValue string
}

func (t *TUI) StartTUI() {
	app := tview.NewApplication()
	// 初始化服务器状态信息

	t.serverInfoView = tview.NewTextView()
	t.serverInfoView.SetBorder(true)
	t.serverInfoView.SetTitle("Redis Server Info")
	t.refreshServerInfo()

	// 初始化历史记录
	t.cmdHistoryView = tview.NewTextView()
	t.cmdHistoryView.SetBorder(true)
	t.cmdHistoryView.SetTitle("CMD History")

	// 初始化命令输入框
	t.cmdInputView = tview.NewInputField().SetLabel(">")
	t.cmdInputView.SetBorder(true)
	t.cmdInputView.SetTitle("cmd")
	t.cmdInputView.SetDoneFunc(t.inputDoneFunc)

	// 初始化 key 信息界面
	t.keyInfoView = tview.NewTextView()
	t.keyInfoView.SetBorder(true)
	t.keyInfoView.SetTitle("keyInfo")

	// 初始化 Value 界面
	t.keyValueView = tview.NewTextView()
	t.keyValueView.SetBorder(true)
	t.keyValueView.SetTitle("Value")

	// 初始化 HEX Value 界面
	t.keyHexValueView = tview.NewTextView()
	t.keyHexValueView.SetBorder(true)
	t.keyHexValueView.SetTitle("HEX Value")

	t.allWidgets = []tview.Primitive{
		t.cmdInputView,
		t.serverInfoView,
		t.cmdHistoryView,
		t.keyInfoView,
		t.keyValueView,
		t.keyHexValueView,
	}

	// 初始化布局
	t.mainPage = tview.NewFlex()
	left := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(t.serverInfoView,0,3,false).
		AddItem(t.cmdHistoryView,0,6,false).
		AddItem(t.cmdInputView,3,1,true)
	right := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(t.keyInfoView,0,2,false).
		AddItem(t.keyValueView,0,6,false).
		AddItem(t.keyHexValueView,0,2,false)
	t.mainPage.AddItem(left,0,1,true).AddItem(right,0,2,false)


	// tab 切换对焦
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			t.cycFocus()
		}
		return event
	})


	err := app.SetRoot(t.mainPage,true).EnableMouse(false).Run()
	if err != nil {
		return
	}
}

func (t *TUI) cycFocus()  {
	t.focusIndex++
	if t.focusIndex == len(t.allWidgets) {
		t.focusIndex = 0
	}
	for i,w := range t.allWidgets{
		if i==t.focusIndex {
			w.Focus(func(p tview.Primitive) {

			})
			continue
		}
		w.Blur()
	}
}

func (t *TUI) inputDoneFunc(key tcell.Key)  {
	if key == tcell.KeyEnter {

		// 执行命令
		if strings.TrimSpace(t.cmdInputView.GetText()) != ""{
			cmdByte,selectDB,db := cmd.GetCmdByte(t.cmdInputView.GetText())

			result,dump := t.RedisServer.Conn.ExecCMDByte(cmdByte)
			if selectDB && string(result)=="OK"{
				t.RedisServer.DB=db
				t.refreshServerInfo()
			}
			t.keyValueView.SetText(string(result))
			t.keyHexValueView.SetText(dump)
		}
		t.flushHistory()
	}else if key == tcell.KeyUp{
		t.historyIndex++
		// 选择历史命令
		if t.historyIndex > len(t.TUIData.CmdHistory) {
			t.historyIndex = 1
		}

		if len(t.TUIData.CmdHistory)>0 {
			t.cmdInputView.SetText(t.TUIData.CmdHistory[len(t.TUIData.CmdHistory)-t.historyIndex])
		}

	}else if key == tcell.KeyDown{
		t.historyIndex--
		if t.historyIndex < 1 {
			t.historyIndex = 1
		}
		if len(t.TUIData.CmdHistory)>0 {
			t.cmdInputView.SetText(t.TUIData.CmdHistory[len(t.TUIData.CmdHistory)-t.historyIndex])
		}

	}
}

func (t TUI) refreshServerInfo()  {
	serverInfo := fmt.Sprintf("Addr: %s\nPort: %s\ndb: %s",
		t.RedisServer.Addr,
		t.RedisServer.Port,
		t.RedisServer.DB)
	t.serverInfoView.SetText(serverInfo)
}

func (t *TUI) flushHistory()  {
	t.TUIData.CmdHistory = append(t.TUIData.CmdHistory, t.cmdInputView.GetText())
	cmdHistory1 := strings.Builder{}
	for _,v := range t.TUIData.CmdHistory {
		cmdHistory1.WriteString(v)
		cmdHistory1.WriteString("\n")
	}
	t.cmdHistoryView.SetText(cmdHistory1.String())
	t.historyIndex = 0
	t.cmdInputView.SetText("")
}