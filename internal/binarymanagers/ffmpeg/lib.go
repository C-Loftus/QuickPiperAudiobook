// Copyright 2025 Colton Loftus
// SPDX-License-Identifier: AGPL-3.0-only

package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
)

// Retrieves the duration of an MP3 file in milliseconds using ffprobe.
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
