package ebookconvert

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertToText(t *testing.T) {
	epubPath := filepath.Join("testdata", "test.epub")
	inputFile, err := os.Open(epubPath)
	require.NoError(t, err, "failed to open test EPUB file")
	defer inputFile.Close()

	// Call the ConvertToText function
	outputReader, err := ConvertToText(inputFile, "epub")
	require.NoError(t, err, "ConvertToText returned an error")

	// Read the output to verify its content
	outputBytes, err := io.ReadAll(outputReader)
	require.NoError(t, err, "failed to read output text")

	// Assert the output is non-empty
	require.NotEmpty(t, string(outputBytes), "output text should not be empty")

	// read in the test.txt file and compare it to the output
	expectedOutput, err := os.ReadFile(filepath.Join("testdata", "test.txt"))
	require.NoError(t, err, "failed to read expected output")
	require.Contains(t, string(outputBytes), string(expectedOutput), "output text does not match expected output")
}
