package getstring

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle()
)

type Model struct {
	title  string
	input  textinput.Model // Campos del modelo
	cancel bool
	quit   bool
	err    error
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.input.Value() == "" {
				return m, nil
			}
			m.quit = true
			return m, tea.Quit
		case "ctrl+c", "q":
			m.quit = true
			m.cancel = true
			return m, tea.Quit
		}
	}
	var cmdInput tea.Cmd
	m.input, cmdInput = m.input.Update(msg)
	if m.input.Err != nil {
		m.err = m.input.Err
		m.input.SetValue(m.input.Value()[:len(m.input.Value())-1])
	}
	cmds = append(cmds, cmdInput)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.quit {
		return ""
	}
	var err = ""
	if m.err != nil {
		err = m.err.Error() //TODO: Añadir estilo de error
	}
	return fmt.Sprintf("%s\n%s\n%s\n", titleStyle.Render(m.title), m.input.View(), err)
}

func validator(s string) error {
	if strings.Contains(s, " ") {
		return fmt.Errorf("don't include space in name")
	}
	return nil
}

func New(title string, placeholder string) Model {
	var input = textinput.New()
	input.Placeholder = placeholder
	input.Focus()
	input.Width = 30
	input.CharLimit = 30
	input.Validate = validator
	return Model{
		title: title,
		input: input,
	}
}

func GetNameGeneratedKeys() (string, bool) {
	p := tea.NewProgram(New("Ingresa el nombre con el que quieres reconocer las llaves:", "nombre"))
	var err error
	var m tea.Model
	if m, err = p.Run(); err != nil {
		log.Fatal(err.Error())
	}

	return (m.(Model).input.Value()), m.(Model).cancel
}
