package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// Run the root command
// NOTE: There seems to be an issue with this given the fact
// rootCmd is a global variable and cobra keeps global state
// TODO come back to this and find a way to avoid it.
func executeCommand(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	newRootCmd := rootCmd
	newRootCmd.SetOut(buf)
	newRootCmd.SetErr(buf)
	newRootCmd.SetArgs(args)
	err := newRootCmd.Execute()
	return buf.String(), err
}

func requireExistsThenRemove(t *testing.T, path string) {
	require.FileExists(t, path)
	err := os.Remove(path)
	require.NoError(t, err)
}

func TestRootCommand(t *testing.T) {

	const configData = "mp3: true\nchapters: true\n"
	homedir, err := os.UserHomeDir()
	require.NoError(t, err)
	configPath := filepath.Join(homedir, ".config", "QuickPiperAudiobook", "config.yaml")
	_ = os.Remove(configPath) // fine if it doesn't exist
	err = os.WriteFile(configPath, []byte(configData), 0644)
	require.NoError(t, err)
	_, err = executeCommand("testdata/titlepage_and_2_chapters.epub")
	require.NoError(t, err)

	requireExistsThenRemove(t, "titlepage_and_2_chapters.mp3")

	// make sure cli args override config
	_, err = executeCommand("testdata/titlepage_and_2_chapters.epub", "--mp3=false", "--chapters=false")
	require.NoError(t, err)
	require.NoFileExists(t, "titlepage_and_2_chapters.mp3")

	// make sure the wav file was created
	requireExistsThenRemove(t, "titlepage_and_2_chapters.wav")
}
