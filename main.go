package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	// import the local lib package from the current directory
	lib "QuickPiperAudiobook/lib"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/gen2brain/beeep"
)

func getConvertedRawText(inputPath string) (io.Reader, error) {
	// Ensure input file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		// If the file does not exist, check if it's a URL
		if lib.IsUrl(inputPath) {
			// Download the file
			file, err := lib.DownloadFile(inputPath, filepath.Base(inputPath))
			if err != nil {
				return nil, fmt.Errorf("failed to download file: %v", err)
			}

			// Get the absolute path of the downloaded file
			inputPath, err = filepath.Abs(file.Name())
			if err != nil {
				return nil, fmt.Errorf("failed to get absolute path of file: %v", err)
			}
		} else {
			// Return an error if the file does not exist and it's not a URL
			return nil, fmt.Errorf("input file does not exist: %s", inputPath)
		}
	}

	// Create a temporary file for the intermediate output
	tmpFile, err := os.CreateTemp("", "ebook-convert-*.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	c := color.New(color.Bold, color.FgMagenta)
	c.Println("Converting " + filepath.Base(inputPath) + " to the proper intermediary text format...")

	cmd := exec.Command("ebook-convert", inputPath, tmpFile.Name())

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to convert ebook: %s\nOutput: %s", err, string(output))
	}

	return tmpFile, nil
}

func checkPiperInstalled() bool {
	cmd := exec.Command("which", "piper")
	err := cmd.Run()
	if err != nil {
		if _, err := os.Stat("piper"); os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func installPiper() error {
	fmt.Println("Piper is not in your PATH. Do you want to download it for local use with this script? (yes/no)")

	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) != "yes" && strings.ToLower(response) != "y" {
		return fmt.Errorf("piper installation aborted")
	}

	fmt.Println("Downloading piper...")
	resp, err := http.Get("https://github.com/rhasspy/piper/releases/download/v1.2.0/piper_amd64.tar.gz")
	if err != nil {
		return fmt.Errorf("failed to download piper: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download piper: %s", resp.Status)
	}

	file, err := os.Create("piper.tar.gz")
	if err != nil {
		return fmt.Errorf("failed to create piper tarball file: %v", err)
	}
	defer file.Close()
	defer os.Remove(file.Name())

	fmt.Println("Extracting piper...")
	if err := lib.Untar(resp.Body, "."); err != nil {
		return fmt.Errorf("failed to extract piper: %v", err)
	}

	fmt.Println("Piper installed successfully.")
	return nil
}

func checkEbookConvertInstalled() error {
	_, err := exec.LookPath("ebook-convert")
	if err != nil {
		fmt.Errorf("the ebook-convert command was not found in your PATH. Please install it with your package manager")
	}

	return nil
}

func grabModel(modelName string) error {

	modelURL, ok := modelToURL[modelName]
	if !ok {
		return fmt.Errorf("model not found: %s", modelName)
	}
	modelJSONURL := modelURL + ".json"

	// Download the model
	if err := lib.DownloadIfNotExists(modelURL, modelName); err != nil {
		return err
	}

	if err := lib.DownloadIfNotExists(modelJSONURL, modelName+".json"); err != nil {
		return err
	}

	return nil
}

func findModels(dir string) ([]string, error) {
	// Read the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	var result []string
	for _, file := range files {
		// Skip directories
		if file.IsDir() {
			continue
		}

		name := file.Name()

		// Check if the file has a .onnx extension
		if strings.HasSuffix(name, ".onnx") {
			// Construct the path for the associated .json file
			jsonFile := name + ".json"
			jsonFilePath := filepath.Join(dir, jsonFile)

			// Check if the .json file exists
			if _, err := os.Stat(jsonFilePath); err == nil {
				// If the .json file exists, add the .onnx file path to the result
				abs, err := filepath.Abs(name)
				if err != nil {
					return nil, fmt.Errorf("error getting absolute path: %v", err)
				}

				result = append(result, abs)
			}
		}
	}

	return result, nil
}

func runPiper(filename string, model string, text io.Reader) error {

	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %v", err)
	}

	// Define the path to the 'piper' directory
	piperDir := filepath.Join(currentDir, "piper")

	// Define the path to the 'piper' executable within the 'piper' directory
	piperExecutable := filepath.Join(piperDir, "piper")

	slog.Debug("piper executable path: " + piperExecutable)

	// Output name is equal to the filename with .wav instead of the extension
	outputName := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename)) + ".wav"

	abs, err := filepath.Abs(outputName)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	modelAbs, err := filepath.Abs(model)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	// Construct the command
	cmd := exec.Command(piperExecutable, "--model", modelAbs, "--output_file", abs)
	cmd.Dir = piperDir
	cmd.Stdin = text

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	slog.Debug("Running command: " + strings.Join(cmd.Args, " "))

	// Capture output
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run piper command: %v", err)
	}
	if cmd.Stderr != nil {
		return fmt.Errorf("failed to run piper command: %v", cmd.Stderr)
	}

	beeep.Alert("Audiobook created", "Check the terminal for more info", "")

	return nil
}

type CLI struct {
	Input  string `arg:"" help:"Local path or URL to the input file"`
	Output string `help:"Path to save the converted ebook file (optional)."`
	Model  string `help:"Model to use. (optional)"`
}

func main() {

	var cli CLI
	kong.Parse(&cli)
	m := model{}

	func() {
		if cli.Output == "" {
			cli.Output = "."
		}
		if cli.Model == "" {
			defaultModel := "en_US-hfc_male-medium.onnx"
			cli.Model = defaultModel
			fmt.Printf("No model specified. Defaulting to %s\n", defaultModel)
		}

		if (filepath.Ext(cli.Input)) != ".txt" {

			if err := checkEbookConvertInstalled(); err != nil {
				m.err = err
				return
			}
		}

		// Check if piper is installed and prompt to install if not
		if !checkPiperInstalled() {
			if err := installPiper(); err != nil {
				m.err = err
				return
			}
		}

		models, err := findModels(".")
		if err != nil {
			m.err = err
			return
		}

		if len(models) == 0 {
			fmt.Println("No models found locally")
		} else {
			fmt.Println("Local models found: [ " + strings.TrimSpace(strings.Join(models, " , ")) + " ]")
		}
	}()

	m.spinner = spinner.New()
	m.spinner.Style = spinnerStyle
	m.spinner.Spinner = spinner.Line

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
