package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	activeStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFFF"))
	noStyle = lipgloss.NewStyle()
)

type rootModel struct {
	screenWidth  int
	paneSelected focusIndex
	modelsMap    map[focusIndex]tea.Model
	answers      *Survey
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
	mainModel := NewGreetingModel("Welcome to the Greeter Builder\n")
	previewModel := NewPreviewModel(&srv)
	modelsMap[Main] = mainModel
	modelsMap[Preview] = previewModel

	return &rootModel{
		paneSelected: Main,
		answers:      &srv,
		modelsMap:    modelsMap,
		screenWidth:  maxWidth,
	}
}

func (r *rootModel) Init() tea.Cmd {
	return nil
}

func (r *rootModel) View() string {
	windowSize := r.screenWidth / 2
	var leftView, rightView string
	leftView = r.modelsMap[Main].View()
	rightView = r.modelsMap[Preview].View()

	if r.paneSelected == Main {
		leftView = activeStyle.Width(windowSize).Render(leftView)
	} else {
		leftView = noStyle.Width(windowSize).Render(leftView)
	}

	if r.paneSelected == Preview {
		rightView = activeStyle.Width(windowSize).Render(rightView)
	} else {
		rightView = noStyle.Width(windowSize).Render(rightView)
	}

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
		case "tab":
			r.paneSelected = (r.paneSelected + 1) % 2
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

func ClearScreen() {
	fmt.Print("\033[2J\033[H") //Hack to reset view.
}
