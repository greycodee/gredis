package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/greycodee/gredis/views"
	"os"
)

func main() {
	p1 := tea.NewProgram(views.InitialModel())
	if err := p1.Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
