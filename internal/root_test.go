package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuickPiperAudiobook(t *testing.T) {

	t.Run("E2E", func(t *testing.T) {

		file, err := os.CreateTemp("", "*-test.txt")
		require.NoError(t, err)
		defer file.Close()
		_, err = file.WriteString("This is some test data that will be converted to speech.")
		require.NoError(t, err)

		outputFilename, err := QuickPiperAudiobook(file.Name(), "en_US-lessac-medium.onnx", ".", false, true)
		require.NoError(t, err)
		_, err = os.Stat(outputFilename)
		require.NoError(t, err)
		err = os.Remove(outputFilename)
		require.NoError(t, err)

	})

}
