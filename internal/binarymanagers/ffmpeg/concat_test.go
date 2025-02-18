package ffmpeg

import (
	"QuickPiperAudiobook/internal/binarymanagers"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConcat(t *testing.T) {

	files := []string{"testdata/cow-bell.mp3", "testdata/rooster.mp3"}

	const outputFile = "test_ffmpeg_concat.mp3"
	err := ConcatMp3s(files, outputFile)
	defer os.Remove(outputFile)
	require.NoError(t, err, fmt.Errorf("mp3 output failed to concat with error: %v", err))
	require.FileExists(t, outputFile)
	validateMp3 := exec.Command("ffmpeg", "-v", "error", "-i", outputFile, "-f", "null", "-")
	err = validateMp3.Run()
	require.NoError(t, err, fmt.Errorf("mp3 output validation failed after concat with error: %v", err))

	// make sure the file is bigger than its parts
	concatInfo, err := os.Stat(outputFile)
	require.NoError(t, err)

	firstInfo, err := os.Stat(files[0])
	require.NoError(t, err)
	require.Greater(t, concatInfo.Size(), firstInfo.Size())

	showChapterCmd := fmt.Sprintf("ffprobe -i %s -show_chapters", outputFile)
	output, err := binarymanagers.Run(showChapterCmd)
	require.NoError(t, err)
	chapter1Index := strings.Index(output, "Chapter 1")
	require.Greater(t, chapter1Index, 0)
	chapter2Index := strings.Index(output, "Chapter 2")
	require.Greater(t, chapter2Index, 0)
	require.Greater(t, chapter2Index, chapter1Index)

}
