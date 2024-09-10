package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"QuickPiperAudiobook/lib"

	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"
	"github.com/fatih/color"
)

type CLI struct {
	Input            string `arg:"" help:"Local path or URL to the input file"`
	Output           string `help:"Directory in which to save the converted ebook file (optional)."`
	Model            string `help:"Model to use. (optional)"`
	RemoveDiacritics bool   `help:"Remove diacritics from the input file. (optional)"`
}

func main() {
	// Parse command-line arguments
	var config CLI

	// Get the current user
	usr, _ := user.Current()

	// Construct the config path in the user's home directory
	configDir := filepath.Join(usr.HomeDir, ".config", "QuickPiperAudiobook")
	configPath := filepath.Join(configDir, "config.yaml")

	// Check if the directory exists; if not, create it
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}

		// create config.yaml
		if _, err := os.Create(configPath); err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		fmt.Println("New blank configuration file created at", configPath)
	} else if err != nil {
		fmt.Println("Error checking file:", err)
		return
	}

	parser, _ := kong.New(&config, kong.Configuration(kongyaml.Loader, configPath))

	for _, name := range []string{"output", "model"} {
		_, err := parser.Parse([]string{name})

		if err != nil {
			fmt.Println("Error parsing the value for", name, "in your config file at:", configPath)
			return
		}
	}

	var cli CLI
	ctx := kong.Parse(&cli, kong.Description("Covert a text file to an audiobook using a managed piper install"))

	if cli.Output == "" && config.Output != "" {
		fmt.Println("No output value specified, default from config file: " + config.Output)
		cli.Output = config.Output
	}
	if cli.Model == "" && config.Model != "" {
		fmt.Println("Using model specified in config file: " + config.Model)
		cli.Model = config.Model
	}

	// Set default output path if not provided
	if cli.Output == "" {
		cli.Output = "."
	}
	if strings.HasPrefix(cli.Output, "~/") {
		cli.Output = filepath.Join(usr.HomeDir, cli.Output[2:])
	}

	if cli.Model == "" {
		defaultModel := "en_US-hfc_male-medium.onnx"
		println("No model specified. Defaulting to " + defaultModel)
		cli.Model = defaultModel
	}

	if (filepath.Ext(cli.Input)) != ".txt" {

		if err := lib.CheckEbookConvertInstalled(); err != nil {
			fmt.Printf("Error: %v\n", err)
			ctx.FatalIfErrorf(err)
			return
		}
	}

	// Check if piper is installed and prompt to install if not
	if !lib.CheckPiperInstalled() {
		if err := lib.InstallPiper(); err != nil {
			ctx.FatalIfErrorf(err)
			return
		}
	}

	models, err := lib.FindModels(".")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		ctx.FatalIfErrorf(err)
		return
	}

	if len(models) == 0 {
		fmt.Println("No models found locally")
	} else {
		fmt.Println("Local models found: [ " + strings.TrimSpace(strings.Join(models, " , ")) + " ]")
	}

	err = lib.GrabModel(cli.Model)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		ctx.FatalIfErrorf(err)
		return
	}

	data, err := lib.GetConvertedRawText(cli.Input)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		ctx.FatalIfErrorf(err)
	} else {
		fmt.Println("Text conversion completed successfully.")
	}

	plainText, err := lib.RemoveDiacritics(data)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		ctx.FatalIfErrorf(err)
		return
	}

	err = lib.RunPiper(cli.Input, cli.Model, plainText, cli.Output)

	if err != nil {
		color.Red("Error: %v", err)
		return
	}
}
