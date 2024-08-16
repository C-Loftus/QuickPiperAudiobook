package main

import (
	"fmt"

	// "strings"

	// "github.com/alecthomas/kong"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	textStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	helpStyle    = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type model struct {
	spinner       spinner.Model
	piperModels   []string
	err           error
	selectedModel string
	config        CLI
}

func grabModelCmd(m model) tea.Cmd {
	return func() tea.Msg {
		err := grabModel(m.config.Model)
		if err != nil {
			m.err = err
		}
		return nil
	}
}

func getConvertedRawTextCmd(m model) tea.Cmd {
	return func() tea.Msg {
		data, err := getConvertedRawText(m.config.Input)
		if err != nil {
			m.err = err
		}
		return data
	}
}

func runPiperCmd(m model) tea.Cmd {

	return func() tea.Msg {
		err := runPiper(m.config.Input, m.config.Model, nil)
		if err != nil {
			m.err = err
		}
		return nil
	}
}

func (m model) Init() tea.Cmd {
	if m.err != nil {
		return tea.Quit
	}

	return tea.Sequence(grabModelCmd(m), getConvertedRawTextCmd(m), runPiperCmd(m))

}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.err != nil {
		return m, tea.Quit
	}
	switch msg := msg.(type) {
	// case string:
	// 	m.err = fmt.Errorf(msg)
	// 	return m, m.spinner.Tick
	case nil:
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		default:
			return m, nil
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.QuitMsg:
		return m, tea.Quit
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() (s string) {

	if m.err != nil {
		return m.err.Error()
	} else {
		s += fmt.Sprintf("\n %s %s\n\n", m.spinner.View(), textStyle("Spinning..."))
	}

	return s
}
