package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

type preview struct {
	survey  *Survey
	content string
}

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()
)

func NewPreviewModel(srv *Survey) *preview {
	content, err := yaml.Marshal(srv)
	if err != nil {
		panic(err)
	}
	return &preview{
		survey:  srv,
		content: string(content),
	}

}

func (p *preview) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return p, tea.Quit
		} else {
			//Update model
			content, err := yaml.Marshal(p.survey)
			if err != nil {
				panic(err)
			}
			p.content = string(content)
		}
	}

	return p, nil
}

func (p *preview) Init() tea.Cmd {
	return nil
}

func (p *preview) View() string {
	return titleStyle.Render(p.content)

}
