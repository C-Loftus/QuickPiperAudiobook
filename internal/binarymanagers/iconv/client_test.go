package iconv

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemoveDiacritics(t *testing.T) {
	textToSimplify := "dhyāna jhāna Bön Shān dōng shěng Qí shān"
	piperInput := strings.NewReader(textToSimplify)
	result, err := RemoveDiacritics(piperInput)
	require.NoError(t, err)
	resultBytes, err := io.ReadAll(result)
	require.NoError(t, err)
	require.Equal(t, "dhyana jhana Bon Shan dong sheng Qi shan", string(resultBytes))
}
