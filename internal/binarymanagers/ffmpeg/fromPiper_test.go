// Copyright 2025 Colton Loftus
// SPDX-License-Identifier: AGPL-3.0-only

package ffmpeg

import (
	"os"
	"strings"
	"testing"

	"github.com/C-Loftus/QuickPiperAudiobook/internal/binarymanagers/piper"

	"github.com/stretchr/testify/require"
)

func TestPiperToMp3(t *testing.T) {

	piperClient, err := piper.NewPiperClient("en_US-lessac-medium.onnx")
	require.NoError(t, err)

	const testData = "This is some test data for ffmpeg integration tests."
	const stream = true
	streamData, _, err := piperClient.Run("test_file_name.txt", strings.NewReader(testData), ".", stream)
	require.NoError(t, err)

	file, err := os.CreateTemp("", "ffmpeg_piper_integrated_test_*.mp3")
	require.NoError(t, err)
	defer os.Remove(file.Name())

	err = OutputToMp3(streamData.Stdout, file.Name())
	require.NoError(t, err)
	require.FileExists(t, file.Name())

}
