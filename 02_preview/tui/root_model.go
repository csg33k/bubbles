package tui

import (
	"github.com/76creates/stickers/flexbox"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var activeStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FFFFFFF"))

type rootModel struct {
	flexBox      *flexbox.HorizontalFlexBox
	paneSelected focusIndex
	modelsMap    map[focusIndex]tea.Model
	answers      *Survey
}

type focusIndex uint

const (
	Main    focusIndex = iota
	Preview focusIndex = iota
)

func NewModel() *rootModel {
	srv := NewSurvey()
	modelsMap := make(map[focusIndex]tea.Model)
	mainModel := NewGreetingModel("Welcome to the Greeter Builder\n")
	previewModel := NewPreviewModel(&srv)
	modelsMap[Main] = mainModel
	modelsMap[Preview] = previewModel
	m := flexbox.NewHorizontal(0, 0)
	columns := []*flexbox.Column{
		m.NewColumn().AddCells(
			flexbox.NewCell(1, 1).SetContent(mainModel.View()),
		),
		m.NewColumn().AddCells(
			flexbox.NewCell(1, 1).SetContent(previewModel.View()),
		),
	}
	m.AddColumns(columns)

	return &rootModel{
		flexBox:      m,
		paneSelected: Main,
		answers:      &srv,
		modelsMap:    modelsMap,
	}
}

func (r *rootModel) Init() tea.Cmd {
	return nil
}

func (r *rootModel) View() string {
	return r.flexBox.Render()
}

func (r *rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.flexBox.SetWidth(msg.Width)
		r.flexBox.SetHeight(msg.Height)
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "j":
			r.answers.Age = 42
		case "tab":
			r.paneSelected = (r.paneSelected + 1) % 2
		case "ctrl+c", "esc":
			return r, tea.Quit
		}
	}

	// Update all models
	for ndx, val := range r.modelsMap {
		_, cmd = val.Update(msg)
		if ndx == r.paneSelected {
			content := activeStyle.Render(val.View())
			r.flexBox.GetColumn(int(ndx)).GetCell(0).SetContent(content)
		} else {
			r.flexBox.GetColumn(int(ndx)).GetCell(0).SetContent(val.View())
		}
		cmds = append(cmds, cmd)
		ndx++
	}
	return r, tea.Batch(cmds...)
}

func (r *rootModel) CurrentModal() tea.Model {
	return r.modelsMap[r.paneSelected]
}
