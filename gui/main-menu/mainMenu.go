package mainMenu

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Italic(true).
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("164"))

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("246"))

	lineTopStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false).BorderTop(true).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"})

	selectedButtonStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#7D56F4")).
				Foreground(lipgloss.Color("255")).
				Width(45).
				Align(lipgloss.Center).
				Margin(1, 0, 0, 0)
	unselectedButtonStyle = selectedButtonStyle.BorderForeground(lipgloss.Color("239"))
)

//

const (
	GenerateKeysSelected = iota
	SignSelected
	VerifySelected
	ExitSelected
)

type Model struct {
	selection int
	exit      bool
	altScreen bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd = make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyDown:
			m.selection++
			if m.selection > 3 {
				m.selection = 0
			}
		case tea.KeyUp:
			m.selection--
			if m.selection < 0 {
				m.selection = 3
			}
		case tea.KeyEnter:
			m.exit = true
			return m, tea.Quit
		case tea.KeyCtrlF:
			m.altScreen = !m.altScreen
			if m.altScreen {
				cmds = append(cmds, tea.EnterAltScreen)
			} else {
				cmds = append(cmds, tea.ExitAltScreen)

			}
		case tea.KeyCtrlC:
			m.exit = true
			m.selection = ExitSelected
			return m, tea.Quit
		}
	}

	return m, tea.Batch(cmds...)
}
func GetTitle() string {
	return titleStyle.Render("Practica 6") + "\n\n" +
		"Algoritmo de Firma Digital con Curvas Elipticas." + "\n" +
		lineTopStyle.Render("Para la materia de "+subtitleStyle.Render("Selected Topics in Cryptography."))
}

func GetHelp() string {
	return helpStyle.Render("Enter: Select ● ↓/↑: Move")
}

func (m Model) View() string {
	var styleButtons [4]lipgloss.Style
	styleButtons[0] = unselectedButtonStyle
	styleButtons[1] = unselectedButtonStyle
	styleButtons[2] = unselectedButtonStyle
	styleButtons[3] = unselectedButtonStyle
	styleButtons[m.selection] = selectedButtonStyle
	if m.exit {
		return ""
	}
	return lipgloss.NewStyle().
		Margin(2).
		Render(GetTitle() + "\n\n" +
			lipgloss.JoinVertical(lipgloss.Center,
				styleButtons[0].Render("Generar llaves"),
				styleButtons[1].Render("Firmar Documento"),
				styleButtons[2].Render("Verificar firma")+"\n",
				styleButtons[3].Render("Salir")+"\n",
			) +
			"\n\n" + GetHelp() + "\n")
}

func New() Model {
	return Model{
		selection: GenerateKeysSelected,
	}
}

// Returns the input selected
func ObtenerOpcion() int {
	p := tea.NewProgram(New())
	var m tea.Model
	var err error
	if m, err = p.Run(); err != nil {
		log.Fatal(err.Error())
	}
	return m.(Model).selection
}
