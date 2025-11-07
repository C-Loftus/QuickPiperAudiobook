package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/C-Loftus/QuickPiperAudiobook/internal/binarymanagers"

	"github.com/stretchr/testify/require"
)

func TestConcatWithTitles(t *testing.T) {

	files := []Mp3Section{{Mp3File: "testdata/cow-bell.mp3", Title: "Cow"}, {Mp3File: "testdata/rooster.mp3", Title: "Rooster"}}

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

	firstInfo, err := os.Stat(files[0].Mp3File)
	require.NoError(t, err)
	require.Greater(t, concatInfo.Size(), firstInfo.Size())

	showChapterCmd := []string{"ffprobe", "-i", outputFile, "-show_chapters"}
	output, err := binarymanagers.Run(showChapterCmd)
	require.NoError(t, err)
	chapter1Index := strings.Index(output, files[0].Title)
	require.Greater(t, chapter1Index, 0)
	chapter2Index := strings.Index(output, files[1].Title)
	require.Greater(t, chapter2Index, 0)
	require.Greater(t, chapter2Index, chapter1Index)

}

// Make sure that if we can't get the chapter titles, we don't crash
// and we just fill in generic chapter names like Chapter 1 / Chapter 2
func TestConcatWithoutTitles(t *testing.T) {

	files := []Mp3Section{{Mp3File: "testdata/cow-bell.mp3"}, {Mp3File: "testdata/rooster.mp3"}}

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

	firstInfo, err := os.Stat(files[0].Mp3File)
	require.NoError(t, err)
	require.Greater(t, concatInfo.Size(), firstInfo.Size())

	showChapterCmd := []string{"ffprobe", "-i", outputFile, "-show_chapters"}
	output, err := binarymanagers.Run(showChapterCmd)
	require.NoError(t, err)
	chapter1Index := strings.Index(output, "Chapter 1")
	require.Greater(t, chapter1Index, 0)
	chapter2Index := strings.Index(output, "Chapter 2")
	require.Greater(t, chapter2Index, 0)
	require.Greater(t, chapter2Index, chapter1Index)

}

func TestConcatNonExistFileFails(t *testing.T) {
	files := []Mp3Section{{Mp3File: "DUMMY.mp3", Title: "DUMMY"}, {Mp3File: "DUMMY.mp3", Title: "DUMMY2"}}
	const outputFile = "test_ffmpeg_concat.mp3"
	err := ConcatMp3s(files, outputFile)
	defer os.Remove(outputFile)
	require.Error(t, err)
}
