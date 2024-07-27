package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
)

func ConvertEbook(inputPath string) (io.Reader, error) {
	// Ensure input file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		// If the file does not exist, check if it's a URL
		if IsUrl(inputPath) {
			// Download the file
			file, err := DownloadFile(inputPath, filepath.Base(inputPath))
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

	// Build the ebook-convert command
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

	if strings.ToLower(response) != "yes" {
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

func checkEbookConvertInstalled() error {
	_, err := exec.LookPath("ebook-convert")
	if err != nil {
		return fmt.Errorf("ebook-convert not found in PATH. Please install it with your package manager")
	}

	return nil

}

func grabModel(modelName string) error {

	modelURL := modelToURL[modelName]
	modelJSONURL := modelURL + ".json"

	// Download the model
	if err := downloadIfNotExists(modelURL, modelName); err != nil {
		return err
	}

	if err := downloadIfNotExists(modelJSONURL, modelName+".json"); err != nil {
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
	fmt.Println("Running piper on " + filename)

	// buf := new(strings.Builder)
	// _, err := io.Copy(buf, text)
	// // check errors
	// fmt.Println(buf.String())

	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %v", err)
	}

	// Define the path to the 'piper' directory
	piperDir := filepath.Join(currentDir, "piper")

	// Define the path to the 'piper' executable within the 'piper' directory
	piperExecutable := filepath.Join(piperDir, "piper")

	fmt.Println("piper executable path: " + piperExecutable)

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

	// Print command details
	print("Running command: "+strings.Join(cmd.Args, " "), "\n")

	// Capture output
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run piper command: %v", err)
	}

	fmt.Println("Done. Saved audiobook as", abs)

	return nil
}

type CLI struct {
	Input  string `arg:"" help:"Local path or URL to the input file"`
	Output string `help:"Path to save the converted ebook file (optional)."`
	Model  string `help:"Model to use. (optional)"`
}

func main() {
	// Parse command-line arguments
	var cli CLI
	ctx := kong.Parse(&cli)

	// Set default output path if not provided
	if cli.Output == "" {
		cli.Output = "."
	}

	if cli.Model == "" {
		cli.Model = "en_US-lessac-medium.onnx"
	}

	if err := checkEbookConvertInstalled(); err != nil {
		fmt.Printf("Error: %v\n", err)
		ctx.FatalIfErrorf(err)
		return
	}

	// Check if piper is installed and prompt to install if not
	if !checkPiperInstalled() {
		if err := installPiper(); err != nil {
			ctx.FatalIfErrorf(err)
			return
		}
	}

	models, err := findModels(".")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		ctx.FatalIfErrorf(err)
		return
	}

	fmt.Println("Models found: " + strings.Join(models, ", "))

	err = grabModel(cli.Model)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		ctx.FatalIfErrorf(err)
		return
	}

	// Execute the ebook conversion
	data, err := ConvertEbook(cli.Input)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		ctx.FatalIfErrorf(err)
	} else {
		fmt.Println("Ebook conversion completed successfully.")
	}

	err = runPiper(cli.Input, cli.Model, data)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		ctx.FatalIfErrorf(err)
	} else {
		fmt.Println("Piper completed successfully.")
	}
}
