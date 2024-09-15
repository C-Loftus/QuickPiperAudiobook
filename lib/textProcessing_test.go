package lib

import (
	"os"
	"testing"
)

func TestRemoveDiacritics_Success(t *testing.T) {
	// Create a temporary input file with diacritic characters
	inputFile, err := os.CreateTemp("", "input-*.txt")
	if err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}
	defer os.Remove(inputFile.Name())

	// Write diacritic text to the input file
	inputContent := "Café résumé naïve"
	if _, err := inputFile.WriteString(inputContent); err != nil {
		t.Fatalf("failed to write to input file: %v", err)
	}
	if err := inputFile.Close(); err != nil {
		t.Fatalf("failed to close input file: %v", err)
	}

	// Call RemoveDiacritics function
	outputFile, err := RemoveDiacritics(inputFile)
	if err != nil {
		t.Fatalf("RemoveDiacritics failed: %v", err)
	}
	defer os.Remove(outputFile.Name())

	// Read content from the output file
	outputContent, err := os.ReadFile(outputFile.Name())
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	// Verify the output content
	expectedContent := "Cafe resume naive"
	if string(outputContent) != expectedContent {
		t.Errorf("unexpected output content: got %q, want %q", string(outputContent), expectedContent)
	}
}

func TestRemoveDiacriticsFromFile_Success(t *testing.T) {
	// read in test_diacritics.txt
	inputFile, err := os.Open("test_diacritics.txt")
	if err != nil {
		t.Fatalf("failed to open input file: %v", err)
	}
	defer inputFile.Close()

	// Call RemoveDiacritics function
	outputFile, err := RemoveDiacritics(inputFile)
	if err != nil {
		t.Fatalf("RemoveDiacritics failed: %v", err)
	}
	defer os.Remove(outputFile.Name())

	// Read content from the output file
	outputContent, err := os.ReadFile(outputFile.Name())
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	// Verify the output content
	expectedContent := "dhyana jhana Bon Shan dong sheng Qi shan"
	if string(outputContent) != expectedContent {
		t.Errorf("unexpected output content: got %q, want %q", string(outputContent), expectedContent)
	}
}

func TestRemoveDiacritics_FileNotExist(t *testing.T) {
	// Create a dummy file path that does not exist
	filePath := "/path/to/nonexistent/file.txt"
	// Open the file, this should return an error because the file does not exist
	_, err := os.Open(filePath)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("expected file not exist error but got %v", err)
	}

	// Call RemoveDiacritics function with a nil file reference (since file doesn't exist)
	_, err = RemoveDiacritics(nil)
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err.Error() != "file is nil" {
		t.Errorf("unexpected error: got %v, want %v", err.Error(), "file is nil")
	}
}
