package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/csg33k/bubbles/02_preview/tui"
)

func main() {
	p := tea.NewProgram(tui.NewModel())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
