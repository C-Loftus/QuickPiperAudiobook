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
	os.Args = os.Args[:1]
	os.Args = append(os.Args, "https://example-files.online-convert.com/document/txt/example.txt", "--speak-diacritics")
	RunCLI()
}
