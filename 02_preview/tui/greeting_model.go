package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type greeting struct {
	greeting string
}

func NewGreetingModel(text string) *greeting {
	return &greeting{
		greeting: text,
	}
}

func (g *greeting) Init() tea.Cmd {
	return nil
}

func (g *greeting) View() string {
	msg := fmt.Sprintf("%s\n\n", g.greeting)
	return msg

}

func (g *greeting) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.(type) {
	case tea.KeyMsg:
		key := msg.(tea.KeyMsg).String()
		switch key {
		case "ctrl+c", "esc":
			return g, tea.Quit
		}
	}
	return g, cmd
}
