package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"io"
	"strings"
)

type prompt uint

const (
	namePrompt prompt = iota
	dobPrompt  prompt = iota
	profession prompt = iota
	done       prompt = iota
)

func (p prompt) String() string {
	switch p {
	case namePrompt:
		return "Name"
	case dobPrompt:
		return "Date of Birth"
	case profession:
		return "Profession"
	case done:
		return "Done"
	default:
		return "Unknown"
	}
}

type surveyModel struct {
	state      prompt
	title      string
	survey     *Survey
	errMsg     error
	save       bool
	in         textinput.Model
	savePrompt list.Model
}

func (g *surveyModel) resetInput(title string, validate func(string) error) {
	g.in.SetValue("")
	g.in.Width = 40
	g.in.CharLimit = 156
	g.in.Focus()
	g.in.Placeholder = title
	if validate != nil {
		g.in.Validate = validate

	}
}

var (
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type listItem string

func (i listItem) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItemVal list.Item) {
	i, ok := listItemVal.(listItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func NewGreetingModel(text string, s *Survey) *surveyModel {
	obj := &surveyModel{
		title:  text,
		survey: s,
		in:     textinput.New(),
	}
	items := []list.Item{listItem("Yes"), listItem("No")}

	obj.savePrompt = list.New(items, itemDelegate{}, 0, 0)
	obj.savePrompt.Title = "Save data?"
	obj.savePrompt.SetShowStatusBar(false)
	obj.savePrompt.SetFilteringEnabled(false)
	obj.savePrompt.SetShowHelp(false)
	obj.savePrompt.SetSize(20, 10)
	obj.savePrompt.Styles.Title = titleStyle
	obj.savePrompt.Styles.PaginationStyle = paginationStyle
	obj.savePrompt.Styles.HelpStyle = helpStyle

	obj.resetInput("Name", func(s string) error {
		if s == "" {
			return fmt.Errorf("name cannot be empty")
		}
		return nil
	})

	return obj
}

func (g *surveyModel) Init() tea.Cmd {
	return nil
}

func (g *surveyModel) View() string {
	msg := fmt.Sprintf("%s\n\n", g.title)
	var input string
	var footer = "q to quit"
	if g.state != done {
		input = fmt.Sprintf("%s: %s\n\n", g.in.Placeholder, g.in.View())
		input = noStyle.Border(lipgloss.NormalBorder()).Render(input)
	} else if g.state == done && g.save != true {
		footer = ""
		input = g.savePrompt.View()
		input = wordwrap.String(input, 40)
	} else if g.state == done && g.save == false {
		footer = "Goodbye!"
	}
	return lipgloss.JoinVertical(lipgloss.Left, msg, input, footer)

}

func (g *surveyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	saveInputValue := func(nextState prompt, saveData func(s string)) {
		err := g.in.Validate(g.in.Value())
		if err != nil {
			g.errMsg = err
		} else {
			//g.survey.Name = g.in.Value()
			saveData(g.in.Value())
			g.state = nextState
			g.resetInput(nextState.String(), nil)
		}
	}

	switch msg.(type) {
	case tea.WindowSizeMsg:
		wmsg := msg.(tea.WindowSizeMsg)
		g.savePrompt.SetWidth(wmsg.Width / 2)
		g.savePrompt.SetHeight(wmsg.Height)
	case tea.KeyMsg:
		switch key := msg.(tea.KeyMsg).String(); key {
		case "ctrl+c", "esc":
			return g, tea.Quit
		case "enter":
			switch g.state {
			case namePrompt:
				saveInputValue(dobPrompt, func(s string) {
					g.survey.Name = s
				})
			case dobPrompt:
				saveInputValue(profession, func(s string) {
					g.survey.DOB = s
				})
			case profession:
				saveInputValue(done, func(s string) {
					g.survey.Profession = s
				})
			case done:
				i, ok := g.savePrompt.SelectedItem().(listItem)
				if ok {
					if i == "Yes" {
						g.save = true
						g.Save()
						return g, tea.Quit
					} else if i == "No" {
						g.state = namePrompt
						g.resetInput(namePrompt.String(), nil)
						newData := NewSurvey()
						g.survey.Name = newData.Name
						g.survey.DOB = newData.DOB
						g.survey.Profession = newData.Profession

					}
				} else {
					fmt.Println("unable to get list item")
					return g, tea.Quit
				}

			default:
				fmt.Println("Invalid state detected")
				return g, tea.Quit
			}
		default:
			if g.state != done {
				g.in, cmd = g.in.Update(msg)
				cmds = append(cmds, cmd)
			} else {
				g.savePrompt, cmd = g.savePrompt.Update(msg)
				cmds = append(cmds, cmd)
			}

			if g.state != done {
				_, cmd = g.in.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}

	return g, tea.Batch(cmds...)
}

func (g *surveyModel) Save() {
	fmt.Println("Saving Data model")
}
