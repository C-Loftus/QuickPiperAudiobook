package iconv

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemoveDiacritics(t *testing.T) {

	for _, test := range []struct {
		input    string
		expected string
	}{
		{"dhyāna jhāna Bön Shān dōng sheng Qí shān", "dhyana jhana Bon Shan dong sheng Qi shan"},
		{"test without diacritics", "test without diacritics"},
		{"1234567890", "1234567890"},
		// if you try to remove diacritics on 汉字 it will be replaced by ??
		{"你好", "??"},
	} {
		t.Run(test.input, func(t *testing.T) {
			piperInput := strings.NewReader(test.input)
			result, err := RemoveDiacritics(piperInput)
			require.NoError(t, err)
			resultBytes, err := io.ReadAll(result)
			require.NoError(t, err)
			require.Equal(t, test.expected, string(resultBytes))
		})
	}

}
