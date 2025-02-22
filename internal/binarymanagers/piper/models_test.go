package piper

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExpandModelPath(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	modelName := "test_model"
	modelPath := filepath.Join(tempDir, modelName)
	modelJSONPath := modelPath + ".json"

	t.Run("both onnx and json files are present", func(t *testing.T) {
		err := os.WriteFile(modelPath, []byte("dummy ONNX model"), 0644)
		require.NoError(t, err)
		err = os.WriteFile(modelJSONPath, []byte("dummy JSON"), 0644)
		require.NoError(t, err)

		result, err := expandModelPath(modelName, tempDir)
		if err != nil || result != modelPath {
			t.Errorf("Expected %s, got %s, error: %v", modelPath, result, err)
		}
	})

	t.Run("missing onnx file", func(t *testing.T) {
		os.Remove(modelJSONPath) // remove the JSON file
		result, err := expandModelPath(modelName, tempDir)
		if err == nil || result != "" {
			t.Errorf("Expected error for missing JSON file, got: %v, result: %s", err, result)
		}
	})

	t.Run("model not found", func(t *testing.T) {
		result, err := expandModelPath("non_existent_model", tempDir)
		if err == nil || result != "" {
			t.Errorf("Expected error for non-existent model, got: %v, result: %s", err, result)
		}
	})

	t.Run("model in default directory", func(t *testing.T) {

		modelNameInDir := "another_model"
		modelPathInDir := filepath.Join(tempDir, modelNameInDir)
		modelJSONPathInDir := modelPathInDir + ".json"

		err := os.WriteFile(modelPathInDir, []byte("dummy ONNX model"), 0644)
		require.NoError(t, err)
		err = os.WriteFile(modelJSONPathInDir, []byte("dummy JSON"), 0644)
		require.NoError(t, err)

		result, err := expandModelPath(modelNameInDir, tempDir)
		if err != nil || result != modelPathInDir {
			t.Errorf("Expected %s, got %s, error: %v", modelPathInDir, result, err)
		}
	})
}

// Make sure that the voices.json file is up to date
func TestVoicesJSONIsUpdated(t *testing.T) {
	file, err := os.ReadFile("voices.json")
	require.NoError(t, err)
	fileAsString := string(file)
	const remoteUrl = "https://huggingface.co/rhasspy/piper-voices/raw/main/voices.json"
	resp, err := http.Get(remoteUrl)
	require.NoError(t, err)
	defer resp.Body.Close()
	bodyData := resp.Body
	bodyAsString, err := io.ReadAll(bodyData)
	require.NoError(t, err)
	require.Equal(t, fileAsString, string(bodyAsString))

}
