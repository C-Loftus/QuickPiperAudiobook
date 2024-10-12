package lib

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandModelPath(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	modelName := "test_model"
	modelPath := filepath.Join(tempDir, modelName)
	modelJSONPath := modelPath + ".json"

	// Test case 1: Both ONNX and JSON files are present
	os.WriteFile(modelPath, []byte("dummy ONNX model"), 0644)
	os.WriteFile(modelJSONPath, []byte("dummy JSON"), 0644)

	result, err := ExpandModelPath(modelName, tempDir)
	if err != nil || result != modelPath {
		t.Errorf("Expected %s, got %s, error: %v", modelPath, result, err)
	}

	// Test case 2: ONNX file is present, but JSON file is missing
	os.Remove(modelJSONPath) // remove the JSON file

	result, err = ExpandModelPath(modelName, tempDir)
	if err == nil || result != "" {
		t.Errorf("Expected error for missing JSON file, got: %v, result: %s", err, result)
	}

	// Test case 3: Model not found
	result, err = ExpandModelPath("non_existent_model", tempDir)
	if err == nil || result != "" {
		t.Errorf("Expected error for non-existent model, got: %v, result: %s", err, result)
	}

	// Test case 4: Model found in the default model directory
	modelNameInDir := "another_model"
	modelPathInDir := filepath.Join(tempDir, modelNameInDir)
	modelJSONPathInDir := modelPathInDir + ".json"

	os.WriteFile(modelPathInDir, []byte("dummy ONNX model"), 0644)
	os.WriteFile(modelJSONPathInDir, []byte("dummy JSON"), 0644)

	result, err = ExpandModelPath(modelNameInDir, tempDir)
	if err != nil || result != modelPathInDir {
		t.Errorf("Expected %s, got %s, error: %v", modelPathInDir, result, err)
	}
}
