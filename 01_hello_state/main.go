package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(
		NewModel("Hello world!", "enter your name"),
	)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
