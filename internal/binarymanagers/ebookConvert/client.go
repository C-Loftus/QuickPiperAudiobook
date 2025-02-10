package ebookconvert

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Convert input data to text using the ebook-convert command
// Assumings that the input data is in the format of the file extension provided
// Will output .txt file since piper doesn't support reading other formats
func ConvertToText(input io.Reader, fileExt string) (io.Reader, error) {

	if _, err := exec.LookPath("ebook-convert"); err != nil {
		return nil, fmt.Errorf("the ebook-convert command was not found in your PATH. Please install it with your package manager")
	}

	// have to create a temporary file since ebook-convert doesn't accept stdin
	tmpInputFile, err := os.CreateTemp("", "ebook-convert-temporary-input-*."+fileExt)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer tmpInputFile.Close()
	tmpOutputFile, err := os.CreateTemp("", "ebook-convert-temporary-output-*.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %v", err)
	}

	// Write the input data to the temporary file
	_, err = io.Copy(tmpInputFile, input)
	if err != nil {
		return nil, fmt.Errorf("failed to write to temporary file: %v", err)
	}

	cmd := exec.Command("ebook-convert", tmpInputFile.Name(), tmpOutputFile.Name())

	// make sure that tmpInputFile contains some data and is not an empty file
	if _, err := tmpInputFile.Stat(); err != nil {
		return nil, fmt.Errorf("failed to stat temporary file: %v", err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to convert ebook: %s\nOutput: %s", err, string(output))
	}

	// Return the output file as a reader
	outputFile, err := os.Open(tmpOutputFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to open output file: %v", err)
	}

	return outputFile, nil
}
