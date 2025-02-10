package lib

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDownloadFile(t *testing.T) {
	const url = "https://github.com/C-Loftus/QuickPiperAudiobook/blob/master/readme.md"
	const outputName = "readme.md"
	const dir = "."
	file, err := DownloadFile(url, outputName, ".")
	require.NoError(t, err)
	defer file.Close()
	defer os.Remove(file.Name())
	require.FileExists(t, file.Name())
	require.Equal(t, dir+"/"+outputName, file.Name())
}
