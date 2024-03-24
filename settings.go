package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func save(settings Settings) error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get user config directory: %w", err)
	}

	appDir := "dmvp2p"
	appDirPath := filepath.Join(configDir, appDir)

	if err := os.MkdirAll(appDirPath, 0755); err != nil {
		return fmt.Errorf("error creating app directory: %w", err)
	}
	filename := "settings.json"
	filePath := filepath.Join(appDirPath, filename)

	jsonData, err := json.MarshalIndent(settings, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("error closing file: %v\n", closeErr)
		}
	}()

	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("error writing JSON to file: %w", err)
	}

	return nil
}

func load() (error, Settings) {
	var settings Settings

	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get user config directory: %w", err), settings
	}

	appDir := "dmvp2p"
	filename := "settings.json"
	filePath := filepath.Join(configDir, appDir, filename)

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err), settings
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err), settings
	}

	err = json.Unmarshal(byteValue, &settings)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err), settings
	}

	return nil, settings
}
