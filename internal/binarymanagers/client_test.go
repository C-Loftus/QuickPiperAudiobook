package binarymanagers

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunBinary(t *testing.T) {
	output, err := Run([]string{"echo", "hello"})
	require.NoError(t, err)
	require.Equal(t, "hello\n", output)
}

func FuzzRunPipedBinary(f *testing.F) {
	// Test echo to cat
	f.Fuzz(func(t *testing.T, message string) {
		echoCmd, err := RunPiped("echo", []string{message}, nil)
		require.NoError(t, err)

		catCmd, err := RunPiped("cat", []string{""}, echoCmd.Stdout)
		require.NoError(t, err)

		catCmd2, err := RunPiped("cat", []string{""}, catCmd.Stdout)
		require.NoError(t, err)

		// Properly read from the result.Stdout before calling Wait
		var out bytes.Buffer
		_, err = io.Copy(&out, catCmd2.Stdout)
		require.NoError(t, err)

		// We only need to wait on the last command
		require.NoError(t, catCmd2.Handle.Wait())

		// Assert the final output
		require.Equal(t, message+"\n", out.String())
	})
}
