package lib

import (
	"fmt"
	"os"
)

func CreateConfigIfNotExists(configPath string, configDir string, defaultModel string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("error creating config directory: %v", err)
		}

		defaultConfig := []byte(fmt.Sprintf("output: ~/Audiobooks\nmodel: %q\n", defaultModel))
		if err := os.WriteFile(configPath, defaultConfig, 0644); err != nil {
			return fmt.Errorf("error creating config file: %v", err)
		}
		fmt.Println("New default configuration file created at", configPath)
	} else if err != nil {
		return fmt.Errorf("error checking if config file exists: %v", err)
	}
	return nil
}
