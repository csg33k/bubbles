package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type helloModel struct {
	greeting string
}

func NewModel(text string) helloModel {
	return helloModel{greeting: text}
}

func (s helloModel) Init() tea.Cmd { return nil }

func (s helloModel) View() string {
	return fmt.Sprintf("%s\n\nPress Ctrl+C to exit", s.greeting)
}

func (s helloModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		//String version
		//key := msg.(tea.KeyMsg).String()
		//switch key {
		//case "ctrl+c", "esc":
		//	return s, tea.Quit
		//}
		//KeyType.  String version is more versatile, mainly since not every key combination is available.
		key := msg.(tea.KeyMsg)
		switch key.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return s, tea.Quit
		}
	}
	return s, nil
}
