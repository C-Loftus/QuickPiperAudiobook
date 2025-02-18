package piper

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"

	bin "QuickPiperAudiobook/internal/binarymanagers"
	"QuickPiperAudiobook/internal/lib"
)

type PiperClient struct {
	binary string
	model  string
}

// Install the piper binary to the specified path
func installBinary(installationPath string) error {

	log.Info("Installing piper...")

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

	log.Info("Extracting piper...")
	if err := lib.Untar(resp.Body, installationPath); err != nil {
		return fmt.Errorf("failed to extract piper: %v", err)
	}

	log.Info("Piper installed successfully.")
	return nil
}

func NewPiperClient(model string) (*PiperClient, error) {

	homedir, homedirErr := os.UserHomeDir()
	if homedirErr != nil {
		return nil, fmt.Errorf("failed to get user home directory: %v", homedirErr)
	}

	QuickPiperAudiobookDir := filepath.Join(homedir, ".config", "QuickPiperAudiobook")

	piperDir, err := filepath.Abs(filepath.Join(QuickPiperAudiobookDir, "piper"))
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %v", err)
	}

	if err := os.MkdirAll(piperDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory for piper binary: %v", err)
	}

	// Path to the piper executable within that directory
	piperExecutable := filepath.Join(piperDir, "piper")

	// Check if piper is already installed
	if _, err := os.Stat(piperExecutable); err != nil {
		// Not found, install
		if installErr := installBinary(QuickPiperAudiobookDir); installErr != nil {
			return nil, fmt.Errorf("failed to install piper: %v", installErr)
		}
	}

	fullModelPath, err := findOrDownloadModel(model, QuickPiperAudiobookDir)
	if err != nil {
		return nil, fmt.Errorf("failed to expand model path: %v", err)
	}

	return &PiperClient{model: fullModelPath, binary: piperExecutable}, nil
}

// Run calls piper with the given model, using inputData as the text to be spoken.
//
// If streamOutput == true, it returns a PipedOutput so the caller can read raw PCM
// from output.Stdout.
//
// If streamOutput == false, we wait for piper to finish writing a .wav file and
// run the name of the .wav file as the output.
//
// We log all of piper's stderr so that if there's an error, we see it.
func (p PiperClient) Run(filename string, inputData io.Reader, outdir string, streamOutput bool) (bin.PipedOutput, string, error) {

	absOutdir, err := filepath.Abs(outdir)
	if err != nil {
		return bin.PipedOutput{}, "", fmt.Errorf("failed to get absolute path: %v", err)
	}

	if err := os.MkdirAll(absOutdir, 0755); err != nil {
		return bin.PipedOutput{}, "", fmt.Errorf("output directory %s could not be created: %v", absOutdir, err)
	}

	modelAbs, err := filepath.Abs(p.model)
	if err != nil {
		return bin.PipedOutput{}, "", fmt.Errorf("failed to get absolute path for model: %v", err)
	}

	var outFilePath string
	piperArgs := []string{"-m", modelAbs}

	if streamOutput {
		piperArgs = append(piperArgs, "--output_raw")
	} else {
		// Not streaming, produce a .wav on disk
		baseName := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
		outFilePath = filepath.Join(absOutdir, baseName+".wav")
		outFilePath, err = filepath.Abs(outFilePath)
		if err != nil {
			return bin.PipedOutput{}, "", fmt.Errorf("failed to get absolute path: %v", err)
		}
		piperArgs = append(piperArgs, "--output_file", outFilePath)
	}

	log.Debugf("Running %s with args %v", p.binary, piperArgs)

	output, err := bin.RunPiped(p.binary, piperArgs, inputData)
	if err != nil {
		return bin.PipedOutput{}, "", fmt.Errorf("failed to run piper: %v", err)
	}

	//    If streamOutput is true, we must read from piper's stderr so it doesn't block.
	//    We'll do so in a goroutine that just logs lines. We return immediately so the
	//    caller can read from output.Stdout as well.
	if streamOutput {
		go func() {
			// You can either read line-by-line or read it all at once.
			scanner := bufio.NewScanner(output.Stderr)
			for scanner.Scan() {
				log.Debug(scanner.Text())
			}
			if scanErr := scanner.Err(); scanErr != nil {
				log.Warnf("Error reading piper stderr: %v", scanErr)
			}
		}()

		return output, "", nil
	}

	//    If we are *not* streaming, read the entire stderr in a goroutine or inline
	//    *before* we call Wait(), so we see any messages. We can store them to attach to
	//    error messages if piper fails.

	stderrData, readErr := io.ReadAll(output.Stderr)
	if readErr != nil {
		log.Warnf("Failed to read piper stderr: %v", readErr)
	}
	if len(stderrData) > 0 {
		log.Warnf("piper stderr: %s", string(stderrData))
	}

	waitErr := output.Handle.Wait()
	if waitErr != nil {
		// Include piper's stderr in the error message for clarity
		return bin.PipedOutput{}, "", fmt.Errorf("piper failed: %v.\nStderr:\n%s",
			waitErr, string(stderrData))
	}

	log.Infof("Piper output saved to: %s", outFilePath)
	return bin.PipedOutput{}, outFilePath, nil
}
