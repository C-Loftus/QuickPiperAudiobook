package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func downloadIfNotExists(fileURL, fileName string) error {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		if _, err := DownloadFile(fileURL, fileName); err != nil {
			return err
		}
	}
	return nil
}

func DownloadFile(url string, filename string) (*os.File, error) {

	print("Downloading " + filename)

	// Create the file to save the model
	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("error creating file %s: %v", filename, err)
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
		return nil, fmt.Errorf("error saving file %s: %v", filename, err)
	}

	fmt.Println("Downloaded successfully.")

	return file, nil
}
