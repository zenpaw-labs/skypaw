package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zenpaw-labs/skypaw/ui"
)

func main() {
	// go path_utils.AddToPath()
	m := ui.InitialModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}

}
