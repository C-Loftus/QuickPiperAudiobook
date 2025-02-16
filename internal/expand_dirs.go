package internal

import (
	"os"
	"path/filepath"
	"strings"
)

func expandHomeDir(config AudiobookArgs) (AudiobookArgs, error) {
	expandPath := func(path string) (string, error) {
		if strings.HasPrefix(path, "~") {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return path, err
			}
			return filepath.Join(homeDir, strings.TrimPrefix(path, "~")), nil
		}
		return path, nil
	}

	var err error
	config.OutputDirectory, err = expandPath(config.OutputDirectory)
	if err != nil {
		return config, err
	}

	config.FileName, err = expandPath(config.FileName)
	if err != nil {
		return config, err
	}

	config.Model, err = expandPath(config.Model)
	if err != nil {
		return config, err
	}

	return config, nil
}
