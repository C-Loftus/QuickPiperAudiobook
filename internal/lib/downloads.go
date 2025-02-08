package lib

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func DownloadIfNotExists(fileURL, fileName string, outputDir string) error {

	// add the "/" to the end so it the outputdir is a valid path
	if !strings.HasSuffix(outputDir, "/") {
		outputDir += "/"
	}

	outputPath := outputDir + fileName

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		if _, err := DownloadFile(fileURL, fileName, outputDir); err != nil {
			return err
		}
	}
	return nil
}

func DownloadFile(url string, outputName string, outputDir string) (*os.File, error) {

	if !strings.HasSuffix(outputDir, "/") {
		outputDir += "/"
	}

	outputPath := outputDir + outputName

	println("Downloading " + outputName + " to " + outputPath)

	// Create the file to save the model
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

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status code %d from %s", resp.StatusCode, url)
	}

	// Copy the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error saving file %s: %v", outputPath, err)
	}

	fmt.Println("Finished downloading successfully.")

	return file, nil
}
