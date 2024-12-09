package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	activeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFFF"))
	noStyle    = lipgloss.NewStyle()
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
)

type rootModel struct {
	screenWidth int
	modelsMap   map[focusIndex]tea.Model
	answers     *Survey
}

type focusIndex uint

const (
	Main     focusIndex = iota
	Preview  focusIndex = iota
	maxWidth            = 160
)

func NewModel() *rootModel {
	srv := NewSurvey()
	modelsMap := make(map[focusIndex]tea.Model)
	mainModel := NewGreetingModel("Welcome to the Greeter Builder\n", &srv)
	previewModel := NewPreviewModel(&srv)
	modelsMap[Main] = mainModel
	modelsMap[Preview] = previewModel

	return &rootModel{
		answers:     &srv,
		modelsMap:   modelsMap,
		screenWidth: maxWidth,
	}
}

func (r *rootModel) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen, textinput.Blink)
}

func (r *rootModel) View() string {
	windowSize := r.screenWidth / 2
	leftView := activeStyle.Width(windowSize).Render(r.modelsMap[Main].View())
	rightView := noStyle.Width(windowSize).Render(r.modelsMap[Preview].View())

	return lipgloss.JoinHorizontal(lipgloss.Left, leftView, rightView)
}

func (r *rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.screenWidth = msg.Width
		cmds = append(cmds, tea.ClearScreen)
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "ctrl+c", "esc":
			return r, tea.Quit
		}
	}

	//Update all models
	for _, val := range r.modelsMap {
		_, cmd = val.Update(msg)
		cmds = append(cmds, cmd)
	}
	return r, tea.Batch(cmds...)
}
