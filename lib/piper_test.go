package lib

import (
	"os"
	"path/filepath"
	"testing"
)

func TestModels(t *testing.T) {

	if err := DownloadModelIfNotExists("en_US-hfc_male-medium.onnx", "."); err != nil {
		t.Fatalf("error grabbing model: %v", err)
	}

	models, err := FindModels(".")
	if err != nil {
		t.Fatalf("error finding models: %v", err)
	}

	if len(models) == 0 {
		t.Fatalf("no models found")
	}
}

func TestPiperInstalled(t *testing.T) {

	homeDir, _ := os.UserHomeDir()
	path := filepath.Join(homeDir, ".config/QuickPiperAudiobook")

	if !PiperIsInstalled(path) {
		if err := InstallPiper(path); err != nil {
			t.Fatalf("error installing piper: %v", err)
		}
	}

	if !PiperIsInstalled(path) {
		t.Fatalf("piper should be installed")
	}
}

func TestOutput(t *testing.T) {
	homeDir, _ := os.UserHomeDir()
	path := filepath.Join(homeDir, ".config/QuickPiperAudiobook")

	if !PiperIsInstalled(path) {
		if err := InstallPiper(path); err != nil {
			t.Fatalf("error installing piper: %v", err)
		}
	}

	model := "en_US-hfc_male-medium.onnx"
	if err := DownloadModelIfNotExists(model, "."); err != nil {
		t.Fatalf("error grabbing model: %v", err)
	}
	// create an os.file and write to it
	file, err := os.CreateTemp("", "piper-test-*.txt")
	if err != nil {
		t.Fatalf("error creating temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	if _, err := file.Write([]byte("Hello World")); err != nil {
		t.Fatalf("error writing to file: %v", err)
	}

	err = RunPiper(file.Name(), model, file, "/tmp/")

	if err != nil {
		t.Fatalf("error running piper: %v", err)
	}
	newFile := filepath.Join("/tmp/", filepath.Base(file.Name()))

	if _, err := os.Stat(newFile); os.IsNotExist(err) {
		t.Fatalf("output file not created")
	}
	// print the contents of the file
	content, err := os.ReadFile(newFile)
	if err != nil {
		t.Fatalf("error reading file: %v", err)
	}
	t.Log(string(content))
}
