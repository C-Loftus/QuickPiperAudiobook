package lib

import (
	"testing"
)

func TestModels(t *testing.T) {

	if err := GrabModel("en_US-hfc_male-medium.onnx"); err != nil {
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
	if !CheckPiperInstalled() {
		t.Fatalf("piper is not installed")
	}
}
