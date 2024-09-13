package lib

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
)

func RemoveDiacritics(file *os.File) (*os.File, error) {
	if file == nil {
		return nil, fmt.Errorf("file is nil")
	}

	// Create a temporary file for the output
	tmpFile, err := os.CreateTemp("", "diacritics-removed-*.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %v", err)
	}

	pr, pw := io.Pipe()

	cmd := exec.Command("iconv", "-f", "UTF-8", "-t", "ASCII//TRANSLIT//IGNORE", file.Name())

	cmd.Stdout = pw
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %v", err)
	}

	go func() {
		defer pw.Close()
		if _, err := io.Copy(tmpFile, pr); err != nil {
			fmt.Fprintf(os.Stderr, "failed to copy output to temp file: %v", err)
		}
	}()

	// Wait for the command to complete
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("failed to execute command: %v", err)
	}

	return tmpFile, nil
}

func CheckEbookConvertInstalled() error {
	_, err := exec.LookPath("ebook-convert")
	if err != nil {
		return fmt.Errorf("the ebook-convert command was not found in your PATH. Please install it with your package manager")
	}

	return nil

}

func DownloadModelIfNotExists(modelName string, outputDir string) error {

	modelURL, ok := ModelToURL[modelName]
	if !ok {
		return fmt.Errorf("model not found: %s", modelName)
	}
	modelJSONURL := modelURL + ".json"

	// Download the model
	if err := DownloadIfNotExists(modelURL, modelName, outputDir); err != nil {
		return err
	}

	// Download the model JSON
	if err := DownloadIfNotExists(modelJSONURL, modelName+".json", outputDir); err != nil {
		return err
	}

	return nil
}

func GetConvertedRawText(inputPath string) (*os.File, error) {
	// Ensure input file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		// If the file does not exist, check if it's a URL
		if IsUrl(inputPath) {
			// Download the file
			file, err := DownloadFile(inputPath, filepath.Base(inputPath), ".")
			if err != nil {
				return nil, fmt.Errorf("failed to download file: %v", err)
			}
			defer os.Remove(file.Name())

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
	defer tmpFile.Close()

	c := color.New(color.Bold, color.FgMagenta)
	c.Println("Converting " + filepath.Base(inputPath) + " to the proper intermediary text format...")

	cmd := exec.Command("ebook-convert", inputPath, tmpFile.Name())

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to convert ebook: %s\nOutput: %s", err, string(output))
	}

	return tmpFile, nil
}
