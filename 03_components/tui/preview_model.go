package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

type preview struct {
	survey  *Survey
	hidden  bool
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
		key := msg.String()
		switch key {
		case "ctrl+c", "q", "esc":
			return p, tea.Quit
		}
	}

	//Re-Render content
	content, err := yaml.Marshal(p.survey)
	if err != nil {
		panic(err)
	}
	p.content = string(content)

	return p, nil
}

func (p *preview) Init() tea.Cmd {
	return nil
}

func (p *preview) View() string {
	return titleStyle.Render(p.content)
}
