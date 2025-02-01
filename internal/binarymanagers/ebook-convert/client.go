package ebookconvert

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func checkInstalled() error {
	_, err := exec.LookPath("ebook-convert")
	if err != nil {
		return fmt.Errorf("the ebook-convert command was not found in your PATH. Please install it with your package manager")
	}

	return nil
}

func SplitEpub(input io.Reader) (io.Reader, error) {
	return ConvertToText(input, "epub")
}

func ConvertToText(input io.Reader, fileExt string) (io.Reader, error) {

	if err := checkInstalled(); err != nil {
		return nil, err
	}

	tmpInputFile, err := os.CreateTemp("", "ebook-convert-temporary-input-*."+fileExt)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer tmpInputFile.Close()
	defer os.Remove(tmpInputFile.Name())
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

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to convert ebook: %s\nOutput: %s", err, string(output))
	}

	return tmpOutputFile, nil
}
