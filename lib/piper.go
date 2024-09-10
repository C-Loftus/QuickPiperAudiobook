package lib

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/gen2brain/beeep"
)

func RunPiper(filename string, model string, file *os.File, outdir string) error {

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

	outputName := filepath.Join(outdir, strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))) + ".wav"

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
	cmd.Stdin = file

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	slog.Debug("Running command: " + strings.Join(cmd.Args, " "))

	c := color.New(color.Bold, color.FgMagenta, color.BlinkRapid)

	c.Println("Converting your file to an audiobook...", "This may take a while!")
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner
	s.Start()                                                   // Start the spinner

	// Capture output
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run piper command: %v", err)
	}
	s.Stop()

	color.New(color.Bold, color.FgGreen).Println("Done. Saved audiobook as " + abs)

	beeep.Alert("Audiobook created", "Check the terminal for more info", "")

	return nil
}

func FindModels(dir string) ([]string, error) {
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

func CheckPiperInstalled() bool {
	cmd := exec.Command("which", "piper")
	err := cmd.Run()
	if err != nil {
		if _, err := os.Stat("piper"); os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func InstallPiper() error {
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
	if err := Untar(resp.Body, "."); err != nil {
		return fmt.Errorf("failed to extract piper: %v", err)
	}

	fmt.Println("Piper installed successfully.")
	return nil
}
