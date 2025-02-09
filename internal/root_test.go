package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuickPiperAudiobook(t *testing.T) {

	t.Run("E2E", func(t *testing.T) {
		err := QuickPiperAudiobook("test.txt", "test", "test", true, true)
		require.NoError(t, err)
	})

}
