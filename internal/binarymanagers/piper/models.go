package piper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/c-loftus/QuickPiperAudiobook/internal/lib"
)

// Piper has hundreds of pretrained models on the sample Website
// These are some of the best ones for English. However, as long
// as you have both the .onnx and .onnx.json files locally, you
// can use any model you want or even train your own.
var ModelToURL = map[string]string{
	"en_US-hfc_male-medium.onnx":              "https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/hfc_male/medium/en_US-hfc_male-medium.onnx",
	"en_US-hfc_female-medium.onnx":            "https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/hfc_female/medium/en_US-hfc_female-medium.onnx",
	"en_US-lessac-medium.onnx":                "https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/lessac/medium/en_US-lessac-medium.onnx",
	"en_GB-northern_english_male-medium.onnx": "https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_GB/northern_english_male/medium/en_GB-northern_english_male-medium.onnx",
	// Below is an example of a non-English model
	// I happily accept PRs for others here. It is just a bit tedious to enumerate them all
	// since some do not follow the same pattern.
	"zh_CN-huayan-medium.onnx": "https://huggingface.co/rhasspy/piper-voices/resolve/main/zh/zh_CN/huayan/medium/zh_CN-huayan-medium.onnx",
}

// Try to find the model if it exists and otherwise try to download it
// Return the full path to the model
func findOrDownloadModel(modelName, defaultModelDir string) (string, error) {

	fullModelPath, err := expandModelPath(modelName, defaultModelDir)
	if err == nil {
		return fullModelPath, nil
	}

	modelURL, ok := ModelToURL[modelName]
	if !ok {
		return "", fmt.Errorf("model '%s' not found", modelName)
	}

	file, err := lib.DownloadFile(modelURL, modelName, defaultModelDir)
	if err != nil {
		return "", fmt.Errorf("error downloading model '%s': %v", modelName, err)
	}
	jsonURL := modelURL + ".json"
	_, err = lib.DownloadFile(jsonURL, modelName+".json", defaultModelDir)
	if err != nil {
		return "", fmt.Errorf("error downloading model '%s': %v", modelName, err)
	}
	defer file.Close()
	return file.Name(), nil
}

func expandModelPath(modelName string, defaultModelDir string) (string, error) {
	// when given a modelName check if it is present relatively or in the modelDir
	// a path should only be valid if both the onnx and onnx.json file is present

	if _, err := os.Stat(modelName); err == nil {
		if _, err := os.Stat(modelName + ".json"); err == nil {
			return modelName, nil
		}
		return "", fmt.Errorf("onnx for model '%s' was found but the corresponding onnx.json was not", modelName)
	}

	if _, err := os.Stat(filepath.Join(defaultModelDir, modelName)); err == nil {
		if _, err := os.Stat(filepath.Join(defaultModelDir, modelName) + ".json"); err == nil {
			return filepath.Join(defaultModelDir, modelName), nil
		}
		return "", fmt.Errorf("onnx for model '%s' was found in the model directory: '%s' but the corresponding onnx.json was not", modelName, defaultModelDir)

	}
	return "", fmt.Errorf("model '%s' was not found in the current directory or the default model directory: '%s'", modelName, defaultModelDir)
}

func FindModels(dir string) ([]string, error) {

	if strings.HasPrefix(dir, "~/") {
		usr, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("error getting user home directory: %v", err)
		}
		dir = filepath.Join(usr, dir[2:])
	}

	// Read the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	var models []string
	for _, file := range files {
		// Skip directories
		if file.IsDir() {
			continue
		}

		name := file.Name()

		// Check if the file has a .onnx extension
		if strings.HasSuffix(name, ".onnx") {
			// Construct the path for the associated .json file
			jsonFile := name + ".json"
			jsonFilePath := filepath.Join(dir, jsonFile)

			// Check if the .json file exists
			if _, err := os.Stat(jsonFilePath); err == nil {
				// If the .json file exists, add the .onnx file path to the result
				abs, err := filepath.Abs(name)
				if err != nil {
					return nil, fmt.Errorf("error getting absolute path: %v", err)
				}

				models = append(models, abs)
			}
		}
	}

	return models, nil
}
