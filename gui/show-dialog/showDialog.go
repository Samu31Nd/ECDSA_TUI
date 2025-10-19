package showdialog

import (
	"fmt"
	"log"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	ticksStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
)

type (
	tickMsg struct{}
)

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

type Model struct {
	message  string
	typeErr  bool // TODO: si es tipo error, haz el message con estilo rojo
	Ticks    int
	Quitting bool
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc", "enter":
			m.Quitting = true
			return m, tea.Quit
		}
	case tickMsg:
		if m.Ticks == 0 {
			m.Quitting = true
			return m, tea.Quit
		}
		m.Ticks--
		return m, tick()
	}

	return m, nil
}

func (m Model) View() string {
	if m.Quitting {
		return ""
	}

	return fmt.Sprintf("%s\n\nQuitting in %s seconds", m.message, ticksStyle.Render(strconv.Itoa(m.Ticks)))
}

func New(message string, ticks int) Model {
	return Model{
		message: message,
		Ticks:   ticks,
	}
}

func ShowDialog(message string, timer int) {
	p := tea.NewProgram(New(message, timer))
	if _, err := p.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

func ShowError(message string) {
	p := tea.NewProgram(New("Error: "+message, 5))
	if _, err := p.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
