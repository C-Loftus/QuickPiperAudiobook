package piper

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	bin "QuickPiperAudiobook/internal/binarymanagers"
	"QuickPiperAudiobook/internal/lib"
)

type PiperClient struct {
	binary string
	model  string
}

// Install the piper binary to the specified path
func installBinary(installationPath string) error {

	log.Println("Installing piper...")

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

	log.Println("Extracting piper...")
	if err := lib.Untar(resp.Body, installationPath); err != nil {
		return fmt.Errorf("failed to extract piper: %v", err)
	}

	log.Println("Piper installed successfully.")
	return nil
}

func NewPiperClient(model string) (*PiperClient, error) {

	homedir, homedir_err := os.UserHomeDir()

	if homedir_err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %v", homedir_err)
	}

	QuickPiperAudiobookDir := filepath.Join(homedir, ".config", "QuickPiperAudiobook")

	piperDir, err := filepath.Abs(filepath.Join(QuickPiperAudiobookDir, "piper"))
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %v", err)
	}

	if err := os.MkdirAll(piperDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory for piper binary: %v", err)
	}
	// Define the path to the 'piper' executable within the 'piper' directory
	piperExecutable := filepath.Join(piperDir, "piper")
	if _, err := os.Stat(piperExecutable); err == nil {
		// piper is already installed
	} else {
		if err := installBinary(QuickPiperAudiobookDir); err != nil {
			return nil, fmt.Errorf("failed to install piper: %v", err)
		}
	}

	fullModelPath, err := findOrDownloadModel(model, QuickPiperAudiobookDir)
	if err != nil {
		return nil, fmt.Errorf("failed to expand model path: %v", err)
	}

	return &PiperClient{model: fullModelPath, binary: piperExecutable}, nil
}

// filename must be specified since the file passed in is a tmp file and a dummy name
// file with text to convert
func (piper PiperClient) Run(filename string, inputData io.Reader, outdir string, streamOutput bool) (bin.PipedOutput, string, error) {

	outdir, err := filepath.Abs(outdir)

	if err != nil {
		return bin.PipedOutput{}, "", fmt.Errorf("failed to get absolute path: %v", err)
	}

	// make sure the output directory exists
	err = os.MkdirAll(outdir, 0755)
	if err != nil {
		return bin.PipedOutput{}, "", fmt.Errorf("output directory specified for piper could not be created: %v", err)
	}

	outputName := filepath.Join(outdir, strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))) + ".wav"

	filepathAbs, err := filepath.Abs(outputName)
	if err != nil {
		return bin.PipedOutput{}, "", fmt.Errorf("failed to get absolute path: %v", err)
	}

	modelAbs, err := filepath.Abs(piper.model)
	if err != nil {
		return bin.PipedOutput{}, "", fmt.Errorf("failed to get absolute path: %v", err)
	}

	piperCmd := piper.binary + " -m " + modelAbs
	if streamOutput {
		piperCmd += " --output_raw"
	} else {
		piperCmd += " --output_file " + filepathAbs
	}

	output, err := bin.RunPiped(piperCmd, inputData)
	if err != nil {
		return bin.PipedOutput{}, "", fmt.Errorf("failed to run piper: %v", err)
	}

	if streamOutput {
		return output, "", nil
	} else {
		err = output.Handle.Wait()
		if err != nil {
			return bin.PipedOutput{}, "", fmt.Errorf("failed to wait for piper: %v", err)
		}
		log.Println("Piper output saved to: " + filepathAbs)
		return bin.PipedOutput{}, filepathAbs, nil
	}
}
