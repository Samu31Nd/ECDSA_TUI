package filepick

import (
	"errors"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	fileDescriptionStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("79"))
)

type Model struct {
	Filepicker      filepicker.Model
	SelectedFile    string
	fileDescription string
	quitting        bool
	key             bool
	err             error
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m Model) Init() tea.Cmd {
	return m.Filepicker.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	case clearErrorMsg:
		m.err = nil
	}

	var cmd tea.Cmd
	m.Filepicker, cmd = m.Filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.Filepicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		m.SelectedFile = path
	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.Filepicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		m.SelectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.Filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if m.SelectedFile == "" {

		s.WriteString("Selecciona el archivo ")
		if m.key {
			s.WriteString("que contiene tu ")
		}
		s.WriteString(fileDescriptionStyle.Render(m.fileDescription) + ":")
	} else {
		s.WriteString(fileDescriptionStyle.Render(m.fileDescription))
		if m.key {
			s.WriteString(" seleccionada: ")
		} else {
			s.WriteString(" seleccionado: ")
		}
		s.WriteString(m.Filepicker.Styles.Selected.Render(m.SelectedFile))
	}
	s.WriteString("\n\n" + m.Filepicker.View() + "\n")
	return s.String()
}

func GetKeyFile(fileDescription string) (string, bool) {
	fp := filepicker.New()
	fp.CurrentDirectory = "."
	fp.AllowedTypes = []string{".pem"}

	m := Model{
		Filepicker:      fp,
		fileDescription: fileDescription,
		key:             true,
	}
	tm, _ := tea.NewProgram(&m).Run()
	var cancel bool
	mm := tm.(Model)
	if mm.SelectedFile == "" {
		cancel = true
	}
	return mm.SelectedFile, cancel
}

func GetFile(fileDescription string) (string, bool) {
	fp := filepicker.New()
	fp.CurrentDirectory = "."
	fp.AllowedTypes = []string{".txt", ".md"}

	m := Model{
		Filepicker:      fp,
		fileDescription: fileDescription,
		key:             false,
	}
	tm, _ := tea.NewProgram(&m).Run()
	var cancel bool
	mm := tm.(Model)
	if mm.SelectedFile == "" {
		cancel = true
	}
	return mm.SelectedFile, cancel
}

func GetSignFile(fileDescription string) (string, bool) {
	fp := filepicker.New()
	fp.CurrentDirectory = "."
	fp.AllowedTypes = []string{".signed"}

	m := Model{
		Filepicker:      fp,
		fileDescription: fileDescription,
		key:             false,
	}
	tm, _ := tea.NewProgram(&m).Run()
	var cancel bool
	mm := tm.(Model)
	if mm.SelectedFile == "" {
		cancel = true
	}
	return mm.SelectedFile, cancel
}
