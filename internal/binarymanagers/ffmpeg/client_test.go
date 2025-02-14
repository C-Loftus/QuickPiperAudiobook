package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConcat(t *testing.T) {

	files := []string{"testdata/cow-bell.mp3", "testdata/rooster.mp3"}

	const filename = "test_ffmpeg_concat.mp3"
	err := ConcatMp3s(files, filename)
	defer os.Remove(filename)
	require.NoError(t, err, fmt.Errorf("mp3 output failed to concat %v", err))
	require.FileExists(t, filename)
	validateMp3 := exec.Command("ffmpeg", "-v", "error", "-i", filename, "-f", "null", "-")
	err = validateMp3.Run()
	require.NoError(t, err, fmt.Errorf("mp3 output validation failed after concat %v", err))

	// make sure the file is bigger than its parts
	concatInfo, err := os.Stat(filename)
	require.NoError(t, err)

	firstInfo, err := os.Stat(files[0])
	require.NoError(t, err)
	require.Greater(t, concatInfo.Size(), firstInfo.Size())

}
