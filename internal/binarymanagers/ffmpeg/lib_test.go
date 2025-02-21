package ffmpeg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMp3Duration(t *testing.T) {
	duration, err := getMp3Duration("testdata/cow-bell.mp3")
	require.NoError(t, err)
	require.Equal(t, duration, int64(2115))
}
