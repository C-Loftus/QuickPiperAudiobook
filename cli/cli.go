package cli

import (
	"QuickPiperAudiobook/lib"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"
	"github.com/fatih/color"
)

type CLI struct {
	Input           string `arg:"" help:"Local path or URL to the input file"`
	Output          string `help:"Directory in which to save the converted ebook file"`
	Model           string `help:"Local path to the onnx model for piper to use"`
	SpeakDiacritics bool   `help:"Speak diacritics from the input file"`
	ListModels      bool   `help:"List available models"`
}

// package level variables we want to expose for testing
var homedir, homedir_err = os.UserHomeDir()
var configDir = filepath.Join(homedir, ".config", "QuickPiperAudiobook")
var configFile = filepath.Join(configDir, "config.yaml")

const defaultModel = "en_US-hfc_male-medium.onnx"

// All cli code is outside of the main package for testing purposes
func RunCLI() {

	if homedir_err != nil {
		fmt.Printf("Error getting user home directory: %v\n", homedir_err)
		return
	}

	var config CLI

	if err := lib.CreateConfigIfNotExists(configFile, configDir, defaultModel); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	parser, _ := kong.New(&config, kong.Configuration(kongyaml.Loader, configFile))

	for _, name := range []string{"output", "model"} {
		_, err := parser.Parse([]string{name})

		if err != nil {
			fmt.Println("Error parsing the value for", name, "in your config file at:", configFile)
			return
		}
	}

	var cli CLI
	ctx := kong.Parse(&cli, kong.Description("Covert a text file to an audiobook using a managed piper install"))

	if cli.ListModels {
		models, err := lib.FindModels(configDir)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			ctx.FatalIfErrorf(err)
			return
		}

		if len(models) == 0 {
			fmt.Println("No models found in " + configDir)
		} else {
			fmt.Println("Found models:\n" + strings.TrimSpace(strings.Join(models, "\n")))
		}
		return
	}

	if cli.Output == "" && config.Output != "" {
		fmt.Println("No output value specified, default from config file: " + config.Output)
		cli.Output = config.Output
		// if output is not set and config is not set, default to current directory
	} else if cli.Output == "" && config.Output == "" {
		cli.Output = "."
	}

	if cli.Model == "" && config.Model != "" {
		fmt.Println("Using model specified in config file: " + config.Model)
		cli.Model = config.Model
	} else if cli.Model == "" && config.Model == "" {
		println("No model specified. Defaulting to " + defaultModel)
		cli.Model = defaultModel
	}

	if strings.HasPrefix(cli.Output, "~/") {
		// if it starts with ~, then we need to expand it
		cli.Output = filepath.Join(homedir, cli.Output[2:])
	}

	// if cli.Output does not exist as a directory make it
	if _, err := os.Stat(cli.Output); os.IsNotExist(err) {
		err := os.MkdirAll(cli.Output, os.ModePerm)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			ctx.FatalIfErrorf(err)
			return
		}
	}

	if (filepath.Ext(cli.Input)) != ".txt" {
		// if it is not already a .txt file, then we need to convert it to .txt and thus need ebook-convert
		if err := lib.CheckEbookConvertInstalled(); err != nil {
			fmt.Printf("Error: %v\n", err)
			ctx.FatalIfErrorf(err)
			return
		}
	}

	if !lib.PiperIsInstalled(configDir) {
		if err := lib.InstallPiper(configDir); err != nil {
			ctx.FatalIfErrorf(err)
			return
		}
	} else {
		slog.Debug("Piper install detected in " + configDir)
	}

	modelPath, err := lib.ExpandModelPath(cli.Model, configDir)

	if err != nil {
		// if the path can't be expanded, it doesn't exist and we need to download it
		err := lib.DownloadModelIfNotExists(cli.Model, configDir)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			ctx.FatalIfErrorf(err)
			return
		}
		modelPath, err = lib.ExpandModelPath(cli.Model, configDir)

		if err != nil {
			fmt.Printf("Error could not find the model path after downloading it: %v\n", err)
			ctx.FatalIfErrorf(err)
			return
		}
	}

	data, err := lib.GetConvertedRawText(cli.Input)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		ctx.FatalIfErrorf(err)
	} else if data == nil {
		fmt.Println("After converting" + cli.Input + "to txt, no data was generated.")
		return
	} else {
		fmt.Println("Text conversion completed successfully.")
	}

	if !cli.SpeakDiacritics {
		if data, err = lib.RemoveDiacritics(data); err != nil {
			fmt.Printf("Error: %v\n", err)
			ctx.FatalIfErrorf(err)
			return
		}

	}

	if err := lib.RunPiper(cli.Input, modelPath, data, cli.Output); err != nil {
		color.Red("Error: %v", err)
	}
}
