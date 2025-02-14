package internal

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuickPiperAudiobook(t *testing.T) {

	t.Run("end to end with mp3", func(t *testing.T) {

		file, err := os.CreateTemp("", "*-test.txt")
		require.NoError(t, err)
		defer file.Close()
		_, err = file.WriteString("This is some test data that will be converted to speech.")
		require.NoError(t, err)

		conf := AudiobookArgs{
			FileName:        file.Name(),
			Model:           "en_US-lessac-medium.onnx",
			OutputDirectory: ".",
			SpeakDiacritics: false,
			OutputAsMp3:     true,
			Chapters:        false,
		}

		outputFilename, err := QuickPiperAudiobook(conf)
		defer os.Remove(outputFilename)
		require.NoError(t, err)
		_, err = os.Stat(outputFilename)
		require.NoError(t, err)
		require.True(t, strings.HasSuffix(outputFilename, ".mp3"))
	})

	t.Run("end to end with wav", func(t *testing.T) {

		file, err := os.CreateTemp("", "*-test.txt")
		require.NoError(t, err)
		defer file.Close()
		_, err = file.WriteString("This is some test data that will be converted to speech.")
		require.NoError(t, err)

		conf := AudiobookArgs{
			FileName:        file.Name(),
			Model:           "en_US-lessac-medium.onnx",
			OutputDirectory: ".",
			SpeakDiacritics: false,
			OutputAsMp3:     false,
			Chapters:        false,
		}

		outputFilename, err := QuickPiperAudiobook(conf)
		require.NoError(t, err)
		_, err = os.Stat(outputFilename)
		require.NoError(t, err)
		err = os.Remove(outputFilename)
		require.NoError(t, err)
		require.True(t, strings.HasSuffix(outputFilename, ".wav"))
	})

}
