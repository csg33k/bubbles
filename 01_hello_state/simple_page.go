package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var fgStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingTop(2).
	PaddingLeft(4).
	Width(80)

type helloModel struct {
	useStyle bool
	greeting string
}

func NewModel(text, template string) helloModel {
	return helloModel{
		greeting: text,
	}
}

func (s helloModel) Init() tea.Cmd {
	return textinput.Blink
}

func (s helloModel) View() string {
	msg := fmt.Sprintf("%s\n\nPress Ctrl+C or esc  to exit, t to toggle style", s.greeting)
	if s.useStyle {
		return fgStyle.Render(msg)
	}
	return msg

}

func (s helloModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.(type) {
	case tea.KeyMsg:
		key := msg.(tea.KeyMsg).String()
		switch key {
		case "t":
			s.useStyle = !s.useStyle
		case "ctrl+c", "esc":
			return s, tea.Quit
		}
	}
	return s, cmd
}
