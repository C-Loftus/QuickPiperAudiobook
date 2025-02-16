package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ConcatMp3s concatenates MP3 files with proper chapter metadata.
func ConcatMp3s(mp3sInOrder []string, outputName string) error {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found in PATH: %v", err)
	}

	// Create temporary files for concat list and metadata
	concatFile, err := os.CreateTemp("", "concat-*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temp concat file: %v", err)
	}
	defer os.Remove(concatFile.Name())

	metadataFile, err := os.CreateTemp("", "metadata-*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temp metadata file: %v", err)
	}
	defer os.Remove(metadataFile.Name())

	// Write absolute paths to concat file
	var durations []int64
	for _, mp3 := range mp3sInOrder {
		absPath, err := filepath.Abs(mp3)
		if err != nil {
			return fmt.Errorf("failed to get absolute path of %s: %v", mp3, err)
		}

		_, err = concatFile.WriteString(fmt.Sprintf("file '%s'\n", absPath))
		if err != nil {
			return fmt.Errorf("failed to write to concat file: %v", err)
		}

		// Get duration of MP3 file
		duration, err := getMp3Duration(absPath)
		if err != nil {
			return fmt.Errorf("failed to get duration of %s: %v", absPath, err)
		}
		durations = append(durations, duration)
	}
	err = concatFile.Close() // Ensure file is written
	if err != nil {
		return fmt.Errorf("failed to close concat file: %v", err)
	}

	// Generate metadata file with chapter markers
	if err := generateMetadataFile(metadataFile, durations); err != nil {
		return fmt.Errorf("failed to create metadata file: %v", err)
	}

	// Run ffmpeg to concatenate and embed metadata
	cmd := exec.Command(
		"ffmpeg", "-f", "concat", "-safe", "0", "-i", concatFile.Name(),
		"-i", metadataFile.Name(), "-map_metadata", "1", "-id3v2_version", "3", "-acodec", "libmp3lame",
		"-b:a", "192k", "-y", outputName,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %v\n%s", err, output)
	}

	return nil
}

// getMp3Duration retrieves the duration of an MP3 file in milliseconds using ffprobe.
func getMp3Duration(mp3File string) (int64, error) {
	// make sure that the file exists
	if _, err := os.Stat(mp3File); os.IsNotExist(err) {
		return 0, fmt.Errorf("file %s does not exist", mp3File)
	}

	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1", mp3File)
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe error: %v %s", err, output)
	}

	var durationSec float64
	_, err = fmt.Sscanf(string(output), "%f", &durationSec)
	if err != nil {
		return 0, fmt.Errorf("failed to parse ffprobe output: %v", err)
	}
	return int64(durationSec * 1000), nil // Convert to milliseconds
}

// generateMetadataFile writes an ffmetadata file with chapters based on MP3 durations.
func generateMetadataFile(metadataFile *os.File, durations []int64) error {
	_, err := metadataFile.WriteString(";FFMETADATA1\n")
	if err != nil {
		return err
	}

	startTime := int64(0)
	for i, duration := range durations {
		endTime := startTime + duration
		chapter := fmt.Sprintf("\n[CHAPTER]\nTIMEBASE=1/1000\nSTART=%d\nEND=%d\ntitle=Chapter %d\n",
			startTime, endTime, i+1)

		_, err := metadataFile.WriteString(chapter)
		if err != nil {
			return err
		}

		startTime = endTime
	}
	metadataFile.Close() // Ensure file is written
	return nil
}
