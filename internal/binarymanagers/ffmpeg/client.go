package ffmpeg

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// OutputToMp3 converts raw PCM audio to MP3 using ffmpeg
func OutputToMp3(input io.Reader, outputName string) error {
	// Verify ffmpeg is available
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found in PATH: %v", err)
	}

	// Create FFmpeg command
	cmd := exec.Command("ffmpeg",
		"-f", "s16le", // Raw PCM format
		"-ar", "22050", // Sample rate
		"-ac", "1", // Mono
		"-i", "pipe:0", // Input from stdin
		"-acodec", "libmp3lame",
		"-b:a", "128k", // MP3 bitrate
		"-y", outputName, // Output file
	)

	// Set up ffmpeg stdin
	ffmpegIn, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %v", err)
	}

	// Set up ffmpeg stderr for debugging
	cmd.Stderr = os.Stderr

	// Start ffmpeg
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %v", err)
	}

	// Stream PCM data to ffmpeg
	_, err = io.Copy(ffmpegIn, input)
	if err != nil {
		return fmt.Errorf("failed to write PCM data to ffmpeg: %v", err)
	}

	// Close ffmpeg input to signal end of stream
	ffmpegIn.Close()

	// Wait for ffmpeg to finish
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg command failed: %v", err)
	}

	// Validate MP3 output
	validateCmd := exec.Command("ffmpeg", "-v", "error", "-i", outputName, "-f", "null", "-")
	if err := validateCmd.Run(); err != nil {
		return fmt.Errorf("mp3 output validation failed: %v", err)
	}

	return nil
}
