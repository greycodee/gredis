package views

import (
	"encoding/hex"
	"fmt"
	"github.com/gdamore/tcell/v2"
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

	RedisServerInfo RedisServerInfo
	KeyInfo			KeyInfo
	TUIData			TUIData
}

type RedisServerInfo struct {
	addr 	string
	port 	string
	db 		string
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
	serverInfo := fmt.Sprintf("Addr: %s\nPort: %s\ndb: %s",
		t.RedisServerInfo.addr,
		t.RedisServerInfo.port,
		t.RedisServerInfo.db)
	t.serverInfoView = tview.NewTextView()
	t.serverInfoView.SetText(serverInfo)
	t.serverInfoView.SetBorder(true)
	t.serverInfoView.SetTitle("Redis Server Info")

	// 初始化历史记录
	//cmdHistory := strings.Builder{}
	//for _,v := range t.TUIData.CmdHistory {
	//	cmdHistory.WriteString(v)
	//}
	t.cmdHistoryView = tview.NewTextView()
	t.cmdHistoryView.SetBorder(true)
	t.cmdHistoryView.SetTitle("Redis Server Info")
	//t.cmdHistoryView.SetText(cmdHistory.String())

	// 初始化命令输入框
	t.cmdInputView = tview.NewInputField().SetLabel(">")
	t.cmdInputView.SetBorder(true)
	t.cmdInputView.SetTitle("cmd")
	t.cmdInputView.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			t.TUIData.CmdHistory = append(t.TUIData.CmdHistory, t.cmdInputView.GetText())
			cmdHistory1 := strings.Builder{}
			for _,v := range t.TUIData.CmdHistory {
				cmdHistory1.WriteString(v)
				cmdHistory1.WriteString("\n")
			}
			t.cmdHistoryView.SetText(cmdHistory1.String())
			t.cmdInputView.SetText("")
		}
	})

	// 初始化 key 信息
	t.keyInfoView = tview.NewTextView()
	t.keyInfoView.SetBorder(true)
	t.keyInfoView.SetTitle("keyInfo")

	t.keyValueView = tview.NewTextView()
	t.keyValueView.SetBorder(true)
	t.keyValueView.SetTitle("Value")

	t.keyHexValueView = tview.NewTextView()
	t.keyHexValueView.SetBorder(true)
	t.keyHexValueView.SetTitle("HEX Value")


	// 初始化布局
	t.mainPage = tview.NewFlex()
	left := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(t.serverInfoView,0,3,false).
		AddItem(t.cmdHistoryView,0,6,false).
		AddItem(t.cmdInputView,0,1,true)
	right := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(t.keyInfoView,0,2,false).
		AddItem(t.keyValueView,0,4,false).
		AddItem(t.keyHexValueView,0,4,false)
	t.mainPage.AddItem(left,50,1,true).AddItem(right,0,2,false)



	err := app.SetRoot(t.mainPage,true).EnableMouse(false).Run()
	if err != nil {
		return
	}
}

func cmdInputFunc()  {

}

func DrawFlex() *tview.Flex {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)
	textView.SetText(hex.Dump([]byte("asdasc\r\nasda")))
	textView.SetBorder(true)
	textView.SetTitle("Key Value")

	input := tview.NewInputField().SetLabel(">")
	input.SetBorder(true)
	input.SetTitle("cmd")

	flex := tview.NewFlex().
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(tview.NewBox().SetBorder(true).SetTitle("Redis Server Info"),0,3,false).
				AddItem(tview.NewBox().SetBorder(true).SetTitle("Command History"),0,6,false).
				AddItem(input,0,1,true), 50,1,false).

		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Key Info"), 0, 3, false).
			AddItem(textView, 0, 7, false), 0, 2, false)
	return flex
}