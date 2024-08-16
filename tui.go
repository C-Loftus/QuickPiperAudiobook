package main

import (
	"fmt"

	// "strings"

	// "github.com/alecthomas/kong"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	// "github.com/fatih/color"
)

var (
	textStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	helpStyle    = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type statusMsg int

type checkConvertInstalledResult struct {
	result bool
	err    error
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type checkPiperInstalledResult bool

type findModelsResult []string

type installPiperResult error

type logMessage string

type model struct {
	spinner       spinner.Model
	piperModels   []string
	err           error
	selectedModel string
}

func (m model) Start() tea.Msg {
	var cli CLI
	// ctx := kong.Parse(&cli)

	// Set default output path if not provided
	if cli.Output == "" {
		cli.Output = "."
	}
	if cli.Model == "" {
		defaultModel := "en_US-hfc_male-medium.onnx"
		cli.Model = defaultModel
		return "No model specified. Defaulting to " + defaultModel
	}
	return "test"
}

func (m model) Init() tea.Cmd {

	return m.Start

	// if (filepath.Ext(cli.Input)) != ".txt" {

	// 	if err := checkEbookConvertInstalled(); err != nil {
	// 		fmt.Printf("Error: %v\n", err)
	// 		ctx.FatalIfErrorf(err)
	// 		return nil
	// 	}
	// }

	// Check if piper is installed and prompt to install if not
	// if !checkPiperInstalled() {
	// 	if err := installPiper(); err != nil {
	// 		ctx.FatalIfErrorf(err)
	// 		return nil
	// 	}
	// }

	// models, err := findModels(".")
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// 	ctx.FatalIfErrorf(err)
	// 	return nil
	// }

	// if len(models) == 0 {
	// 	fmt.Println("No models found locally")
	// } else {
	// 	fmt.Println("Local models found: [ " + strings.TrimSpace(strings.Join(models, " , ")) + " ]")
	// }

	// err = grabModel(cli.Model)
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// 	ctx.FatalIfErrorf(err)
	// 	return nil
	// }

	// data, err := getConvertedRawText(cli.Input)

	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// 	ctx.FatalIfErrorf(err)
	// } else {
	// 	fmt.Println("Text conversion completed successfully.")
	// }

	// err = runPiper(cli.Input, cli.Model, data)

	// if err != nil {
	// 	color.Red("Error: %v", err)
	// 	return nil
	// }

}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case string:
		return m, nil
	case nil:
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		default:
			return m, nil
		}
	case errMsg:
		m.err = msg
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case checkConvertInstalledResult:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		return m, nil
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() (s string) {
	s += fmt.Sprintf("\n %s %s\n\n", m.spinner.View(), textStyle("Spinning..."))
	return s
}
