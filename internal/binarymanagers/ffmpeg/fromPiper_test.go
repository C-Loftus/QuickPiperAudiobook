package ffmpeg

import (
	"QuickPiperAudiobook/internal/binarymanagers/piper"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPiperToMp3(t *testing.T) {

	piperClient, err := piper.NewPiperClient("en_US-lessac-medium.onnx")
	require.NoError(t, err)

	const testData = "This is some test data for ffmpeg integration tests."
	const stream = true
	streamData, _, err := piperClient.Run("/tmp/piper_test.txt", strings.NewReader(testData), "/tmp", stream)
	require.NoError(t, err)

	const testMp3OutputName = "/tmp/ffmpeg_piper_integrated_test.mp3"
	err = OutputToMp3(streamData.Stdout, testMp3OutputName)
	defer os.Remove(testMp3OutputName)
	require.NoError(t, err)
	require.FileExists(t, "/tmp/test.mp3")

}
