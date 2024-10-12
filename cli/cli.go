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
	Input      string `arg:"" help:"Local path or URL to the input file"`
	Output     string `help:"The directory in which to save the output audiobook"`
	Model      string `help:"Local path to the onnx model for piper to use"`
	SpeakUTF8  bool   `help:"Speak UTF-8 characters; Necessary for many non-English languages."`
	ListModels bool   `help:"List piper models which are installed locally"`
}

// package level variables we want to expose for testing
var homedir, homedir_err = os.UserHomeDir()
var configDir = filepath.Join(homedir, ".config", "QuickPiperAudiobook")
var configFile = filepath.Join(configDir, "config.yaml")

const defaultModel = "en_US-hfc_male-medium.onnx"

// All cli code is outside of the main package for testing purposes
func RunCLI() {

	if homedir_err != nil {
		fmt.Fprintf(os.Stderr, "Error getting user home directory: %v\n", homedir_err)
		os.Exit(1)
	}

	var config CLI

	if err := lib.CreateConfigIfNotExists(configFile, configDir, defaultModel); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating default config file: %v\n", err)
		os.Exit(1)
	}

	parser, _ := kong.New(&config, kong.Configuration(kongyaml.Loader, configFile))

	for _, name := range []string{"output", "model"} {
		_, err := parser.Parse([]string{name})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing the value for %s in your config file at: %s\n", name, configFile)
			os.Exit(1)
		}
	}

	var cli CLI
	ctx := kong.Parse(&cli, kong.Description("Covert a text file to an audiobook using a managed piper install"))

	if cli.ListModels {
		models, err := lib.FindModels(configDir)
		if err != nil {
			ctx.FatalIfErrorf(err)
		}

		if len(models) == 0 {
			fmt.Println("No models found in " + configDir)
		} else {
			fmt.Println("Found models:\n" + strings.TrimSpace(strings.Join(models, "\n")))
		}
		return
	}

	if cli.Output == "" && config.Output != "" {
		fmt.Println("No output directory specified, default from config file: " + config.Output)
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
			ctx.FatalIfErrorf(err)
		}
	}

	if (filepath.Ext(cli.Input)) != ".txt" {
		// if it is not already a .txt file, then we need to convert it to .txt and thus need ebook-convert
		if err := lib.CheckEbookConvertInstalled(); err != nil {
			fmt.Printf("Error: %v\n", err)
			ctx.FatalIfErrorf(err)
		}
	}

	if !lib.PiperIsInstalled(configDir) {
		if err := lib.InstallPiper(configDir); err != nil {
			ctx.FatalIfErrorf(err)
		}
	} else {
		slog.Debug("Piper install detected in " + configDir)
	}

	modelPath, modelPathErr := lib.ExpandModelPath(cli.Model, configDir)

	// Some errors above are fine; (we can just download the corresponding model)
	// but others that pertain to having the model but not the corresponding metadata are
	// an error that should be fatal
	if modelPathErr != nil && strings.Contains(modelPathErr.Error(), "but the corresponding") {
		ctx.FatalIfErrorf(modelPathErr)
	}

	if modelPathErr != nil {
		// if the path can't be expanded, it doesn't exist and we need to download it
		err := lib.DownloadModelIfNotExists(cli.Model, configDir)

		if err != nil && modelPathErr != nil {
			fmt.Printf("Error: %v\n", modelPathErr)
		}

		if err != nil {
			ctx.FatalIfErrorf(err)
		}
		modelPath, err = lib.ExpandModelPath(cli.Model, configDir)
		if err != nil {
			ctx.FatalIfErrorf(fmt.Errorf("error could not find the model path after downloading it: %v", err))
		}
	}

	data, err := lib.GetConvertedRawText(cli.Input)

	if err != nil {
		ctx.FatalIfErrorf(err)
	} else if data == nil {
		ctx.FatalIfErrorf(fmt.Errorf("after converting %s to txt, no data was generated", cli.Input))
	}

	if !cli.SpeakUTF8 {
		if data, err = lib.RemoveDiacritics(data); err != nil {
			ctx.FatalIfErrorf(err)
		}

	}

	if err := lib.RunPiper(cli.Input, modelPath, data, cli.Output); err != nil {
		color.Red("Error: %v", err)
	}
}
