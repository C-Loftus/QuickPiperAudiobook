package lib

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	log "github.com/charmbracelet/log"
)

// Return true if the string is a URL
func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// Download a file from a URL and put it in the specified directory with the specified name
func DownloadFile(url, outputName, outputDir string) (*os.File, error) {

	if !strings.HasSuffix(outputDir, "/") {
		outputDir += "/"
	}

	outputPath := outputDir + outputName

	log.Print("Downloading " + outputName + " to " + outputPath)

	file, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("error creating file %s: %v", outputPath, err)
	}
	defer file.Close()

	// Make a GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making GET request to %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status code %d from %s with body %s", resp.StatusCode, url, resp.Body)
	}

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		return nil, fmt.Errorf("error saving file %s: %v", outputPath, err)
	}

	log.Info("Finished downloading successfully.")

	return file, nil
}
