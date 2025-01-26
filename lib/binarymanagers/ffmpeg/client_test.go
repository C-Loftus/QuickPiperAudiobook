package ffmpeg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFFmpegClient(t *testing.T) {
	client, err := NewFFmpegClient()
	require.NotNil(t, client)
	require.NoError(t, err)
}
