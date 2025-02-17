package ffmpeg

import (
	"QuickPiperAudiobook/internal/binarymanagers"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/charmbracelet/log"
)

// Convert raw PCM audio from piper to MP3 using ffmpeg
func OutputToMp3(piperRawAudio io.Reader, outputName string) error {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found in PATH: %v", err)
	}

	if piperRawAudio == nil {
		return fmt.Errorf("nil was passed to ffmpeg mp3 generation")
	}

	args := []string{"-f", "s16le", "-ar", "22050", "-ac", "1", "-i", "pipe:0",
		"-acodec", "libmp3lame", "-b:a", "128k", "-y", outputName}

	output, err := binarymanagers.RunPiped("ffmpeg", args, piperRawAudio)
	if err != nil {
		return err
	}

	log.Debugf("Running ffmpeg to create %s", outputName)

	// Read stderr before waiting
	stderrBytes, _ := io.ReadAll(output.Stderr)
	err = output.Handle.Wait()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %v\nstderr: %s", err, string(stderrBytes))
	}

	// Ensure file was written before verification
	fileInfo, err := os.Stat(outputName)
	if err != nil {
		return fmt.Errorf("output file missing: %v", err)
	}
	if fileInfo.Size() == 0 {
		return fmt.Errorf("output file is empty: %s", outputName)
	}

	// Verify output
	verifyCmd := exec.Command("ffmpeg", "-v", "error", "-i", outputName, "-f", "null", "-")
	verifyOutput, err := verifyCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed to verify output: %v\nstderr: %s", err, string(verifyOutput))
	}

	return nil
}
