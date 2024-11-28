package pages

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Settings struct {
	XmrigPath  string
	P2poolPath string
	Fullscreen bool
}

func NewFullscreenButton(w fyne.Window) *widget.Button {
	return widget.NewButton("Fullscreen", func() {
		// Toggle fullscreen mode
		if w.FullScreen() {
			w.SetFullScreen(false)
		} else {
			w.SetFullScreen(true)
		}
	})
}

func SaveSettings(settings Settings) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, "dmvp2p.json")

	jsonData, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize config to JSON: %w", err)
	}
	err = os.WriteFile(configPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}

	fmt.Printf("Configuration saved to: %s\n", configPath)
	return nil
}

func LoadSettings() (Settings, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Settings{}, fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, "dmvp2p.json")
	fileData, err := os.ReadFile(configPath)
	if err != nil {
		return Settings{}, fmt.Errorf("failed to read configuration file: %w", err)
	}

	var settings Settings
	err = json.Unmarshal(fileData, &settings)
	if err != nil {
		return Settings{}, fmt.Errorf("failed to parse JSON configuration: %w", err)
	}

	fmt.Printf("Configuration loaded from: %s\n", configPath)
	return settings, nil
}
