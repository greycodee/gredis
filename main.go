package main

import (
	"github.com/greycodee/gredis/views"
)

func main() {
	//p1 := tea.NewProgram(views.InitialModel(),tea.WithAltScreen())
	//if err := p1.Start(); err != nil {
	//	fmt.Printf("could not start program: %s\n", err)
	//	os.Exit(1)
	//}

	//block := lipgloss.PlaceHorizontal(10, lipgloss.Center, "asd")
	//
	//fmt.Printf("%s\n", block)

	// Tview
	//app := tview.NewApplication()
	//err := app.SetRoot(views.DrawFlex(),true).EnableMouse(true).Run()
	//if err != nil {
	//	return
	//}
	//fmt.Println(views.KeyType(6))
	//in := &views.KeyInfo{
	//	1,
	//	"123",
	//	views.STRING,
	//}
	tui := views.TUI{

	}
	tui.StartTUI()
}
