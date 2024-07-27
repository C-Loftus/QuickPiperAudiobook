package main

import (
	"testing"
)

func TestModels(t *testing.T) {

	if err := grabModel("en_US-hfc_male-medium.onnx"); err != nil {
		t.Fatalf("error grabbing model: %v", err)
	}

	models, err := findModels(".")
	if err != nil {
		t.Fatalf("error finding models: %v", err)
	}

	if len(models) == 0 {
		t.Fatalf("no models found")
	}
}

func TestPiperInstalled(t *testing.T) {
	if !checkPiperInstalled() {
		t.Fatalf("piper is not installed")
	}
}
