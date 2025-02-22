package piper

import (
	"encoding/json"
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
func TestModelToUrl(t *testing.T) {

	urls, err := modelNameToUrls("ar_JO-kareem-low")
	require.NoError(t, err)
	expectedOnnx := "https://huggingface.co/rhasspy/piper-voices/resolve/main/ar/ar_JO/kareem/low/ar_JO-kareem-low.onnx"
	require.Equal(t, expectedOnnx, urls.OnnxFile)

	expectedJson := "https://huggingface.co/rhasspy/piper-voices/resolve/main/ar/ar_JO/kareem/low/ar_JO-kareem-low.onnx.json"
	require.Equal(t, expectedJson, urls.JsonConfigFile)

	t.Run("test all voices in voices.json", func(t *testing.T) {
		data, err := os.ReadFile("testdata/voices.json")
		require.NoError(t, err)
		var jsonData map[string]interface{}
		err = json.Unmarshal(data, &jsonData)
		require.NoError(t, err)
		for model := range jsonData {
			_, err := modelNameToUrls(model)
			require.NoError(t, err)
		}
	})
}
