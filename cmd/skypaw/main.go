package main

import (
	"github.com/zenpaw-labs/skypaw/ui"
	"github.com/zenpaw-labs/skypaw/utils/path_utils"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	go path_utils.AddToPath()
	m := ui.InitialModel("Okhtyrka")
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}

}
