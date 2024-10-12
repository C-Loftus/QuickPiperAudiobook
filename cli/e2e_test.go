package cli

import (
	"os"
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
	os.Args = append(origArgs[:1], "https://example-files.online-convert.com/document/txt/example.txt", "--speak-diacritics")
	RunCLI()

	// make sure that after running you can run the list models command and it will work
	os.Args = append(origArgs[:1], "https://example-files.online-convert.com/document/txt/example.txt", "--list-models")
	RunCLI()
}
