package lib

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/gen2brain/beeep"
)

func RunPiper(filename string, // we need to have the filename here since the file passed in is a tmp file and a dummy name
	model string, file *os.File, outdir string) error {

	// Debugging: Read file content to check if it's empty
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file passed as input to piper: %v", err)
	}

	// Print file content for debugging purposes
	if len(fileContent) == 0 {
		slog.Debug("File is empty.")
	} else {
		slog.Debug("File content: " + string(fileContent))
	}

	// Reset the file read pointer to the start
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to reset file pointer: %v", err)
	}

	// Get the current working directory
	usr, _ := user.Current()

	piperDir, _ := filepath.Abs(filepath.Join(usr.HomeDir, ".config", "QuickPiperAudiobook", "piper"))

	// Define the path to the 'piper' executable within the 'piper' directory
	piperExecutable := filepath.Join(piperDir, "piper")

	slog.Debug("piper executable path: " + piperExecutable)

	outdir, _ = filepath.Abs(outdir)
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

	if strings.HasPrefix(dir, "~/") {
		usr, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("error getting user home directory: %v", err)
		}
		dir = filepath.Join(usr, dir[2:])
	}

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

func PiperIsInstalled(installationPath string) bool {

	if _, err := os.Stat(installationPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func InstallPiper(installationDir string) error {
	fmt.Println("Installing piper...")

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
	if err := Untar(resp.Body, installationDir); err != nil {
		return fmt.Errorf("failed to extract piper: %v", err)
	}

	fmt.Println("Piper installed successfully.")
	return nil
}
