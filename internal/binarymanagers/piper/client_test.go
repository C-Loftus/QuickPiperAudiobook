package piper

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func cleanupConfigDir(t *testing.T) string {
	homedir, err := os.UserHomeDir()
	require.NoError(t, err)
	QuickPiperAudiobookDir := filepath.Join(homedir, ".config", "QuickPiperAudiobook")
	err = os.RemoveAll(QuickPiperAudiobookDir)
	require.NoError(t, err)
	return QuickPiperAudiobookDir
}

func TestPiperClient(t *testing.T) {

	t.Run("installs binaries", func(t *testing.T) {
		dir := cleanupConfigDir(t)
		client, err := NewPiperClient("en_US-lessac-medium.onnx")
		require.NoError(t, err)
		require.Equal(t, filepath.Join(dir, "piper", "piper"), client.binary)
		_, err = exec.LookPath(client.binary)
		require.NoError(t, err)
	})

	t.Run("converts data", func(t *testing.T) {
		client, err := NewPiperClient("en_US-lessac-medium.onnx")
		require.NoError(t, err)
		_, outputFilename, err := client.Run("test_file_name.txt", strings.NewReader("This is some test data for piper integration tests."), ".", false)
		require.NoError(t, err)
		defer os.Remove(outputFilename)
		require.FileExists(t, outputFilename)
	})

}
