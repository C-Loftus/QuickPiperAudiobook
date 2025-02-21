package lib

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDownloadFile(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Test file content")
	}))
	defer server.Close()

	const outputName = "readme.md"
	const dir = "."
	url := server.URL + "/testfile"

	// Call the function with the local server URL
	file, err := DownloadFile(url, outputName, dir)
	require.NoError(t, err)
	defer file.Close()
	defer os.Remove(file.Name())

	// Ensure file exists
	require.FileExists(t, file.Name())

	// Ensure correct file path
	require.Equal(t, dir+"/"+outputName, file.Name())

	// Check file contents
	content, err := os.ReadFile(file.Name())
	require.NoError(t, err)
	require.Equal(t, "Test file content\n", string(content))
}
