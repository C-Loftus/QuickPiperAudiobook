package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCLI(t *testing.T) {
	// reset all cli args, since the golang testing framework changes them
	os.RemoveAll(configDir)
	os.Args = os.Args[:1]
	os.Args = append(os.Args, "https://example-files.online-convert.com/document/txt/example.txt")
	RunCLI()
}

func TestCLIWithDiacritics(t *testing.T) {
	// reset all cli args, since the golang testing framework changes them
	os.RemoveAll(configDir)
	origArgs := os.Args
	os.Args = append(origArgs[:1], "https://example-files.online-convert.com/document/txt/example.txt", "--speak-utf-8")
	RunCLI()

	// make sure that after running you can run the list models command and it will work
	os.Args = append(origArgs[:1], "https://example-files.online-convert.com/document/txt/example.txt", "--list-models")
	RunCLI()
}

// Test that the cli works with chinese language text
func TestChinese(t *testing.T) {
	// reset all cli args, since the golang testing framework changes them
	os.RemoveAll(configDir)
	origArgs := os.Args
	// get the file located at ../lib/test_chinese.txt
	os.Args = append(origArgs[:1], "../lib/test_chinese.txt", "--model=zh_CN-huayan-medium.onnx", "--speak-utf-8")
	RunCLI()

	// check if there is a file at ~/Audiobooks/test_chinese.wav and make sure it's not empty
	homedir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("error getting user home directory: %v", err)
	}
	testFile := filepath.Join(homedir, "Audiobooks", "test_chinese.wav")
	defer os.Remove(testFile)

	if info, err := os.Stat(testFile); err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("file not created: %v", err)
		}
		t.Fatalf("error getting file info: %v", err)
	} else if info.Size() == 0 {
		t.Fatalf("file is empty")
	}

}
