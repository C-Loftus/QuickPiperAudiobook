package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Mp3Section struct {
	// The title of the chapter
	Title string
	// The path to the mp3 file to use when concatenating
	Mp3File string
	// Duration of the MP3 file in milliseconds
	Duration int64
}

// Concatenates MP3 files and saves the output as an MP3 file
// with proper chapter metadata markers
func ConcatMp3s(sectionsInOrder []Mp3Section, outputName string) error {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found in PATH: %v", err)
	}

	// Create temporary files for concat list and metadata
	// this is needed for ffmpeg since ffmpeg uses it to determine the order of the files
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

	for i, section := range sectionsInOrder {
		absPath, err := filepath.Abs(section.Mp3File)
		if err != nil {
			return fmt.Errorf("failed to get absolute path of %v: %v", section, err)
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
		sectionsInOrder[i].Duration = duration // Update the original slice element since we are iterating by copied value
	}
	err = concatFile.Close() // Ensure file is written
	if err != nil {
		return fmt.Errorf("failed to close concat file: %v", err)
	}

	// Generate metadata file with chapter markers
	if err := generateMetadataFile(metadataFile, sectionsInOrder); err != nil {
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

// Write an ffmetadata file with chapters based on MP3 durations.
func generateMetadataFile(metadataFile *os.File, sectionsInOrder []Mp3Section) error {
	_, err := metadataFile.WriteString(";FFMETADATA1\n")
	if err != nil {
		return err
	}

	startTime := int64(0)
	for i, section := range sectionsInOrder {
		endTime := startTime + section.Duration

		if section.Title == "" {
			section.Title = fmt.Sprintf("Chapter %d", i+1)
		}

		// Start chapters 500ms before the actual start time
		// Without this, the chapter starts exactly when speech starts
		// which often ends up cutting off the first word or making
		// it sound bad in many audiobook players
		const fiveHundredMs = 500
		startTimeOffset := startTime - fiveHundredMs
		if startTimeOffset < 0 {
			startTimeOffset = 0
		}

		chapter := fmt.Sprintf("\n[CHAPTER]\nTIMEBASE=1/1000\nSTART=%d\nEND=%d\ntitle=%s\n",
			startTimeOffset, endTime, section.Title)

		_, err := metadataFile.WriteString(chapter)
		if err != nil {
			return err
		}

		startTime = endTime
	}
	metadataFile.Close()
	return nil
}
