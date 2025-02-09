package piper

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func CleanupConfigDir(t *testing.T) string {
	homedir, err := os.UserHomeDir()
	require.NoError(t, err)
	QuickPiperAudiobookDir := filepath.Join(homedir, ".config", "QuickPiperAudiobook")
	err = os.RemoveAll(QuickPiperAudiobookDir)
	require.NoError(t, err)
	return QuickPiperAudiobookDir
}

func TestPiperClient(t *testing.T) {

	t.Run("installs binaries", func(t *testing.T) {
		dir := CleanupConfigDir(t)
		client, err := NewPiperClient("en_US-lessac-medium.onnx")
		require.NoError(t, err)
		require.Equal(t, filepath.Join(dir, "piper", "piper"), client.binary)
		_, err = exec.LookPath(client.binary)
		require.NoError(t, err)
	})

}
