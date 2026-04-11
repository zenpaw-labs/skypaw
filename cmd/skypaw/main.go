package main

import (
	"skypaw/ui"
	"skypaw/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	go utils.AddToPath()
	m := ui.InitialModel("Okhtyrka")
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
