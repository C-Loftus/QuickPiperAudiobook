package lib

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateConfigIfNotExists(t *testing.T) {
	homedir, err := os.UserHomeDir()

	if err != nil {
		t.Fatalf("error getting user home directory: %v", err)
	}

	configDir := filepath.Join(homedir, ".config", "QuickPiperAudiobook")
	configFile := filepath.Join(configDir, "config.yaml")
	defaultModel := "en_US-hfc_male-medium.onnx"

	if err := CreateConfigIfNotExists(configFile, configDir, defaultModel); err != nil {
		t.Fatalf("error creating config file: %v", err)
	}

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Fatalf("config directory not created: %v", err)
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Fatalf("config file not created: %v", err)
	}

	//teardown
	if err := os.Remove(configFile); err != nil {
		t.Fatalf("error removing config file: %v", err)
	}

}
