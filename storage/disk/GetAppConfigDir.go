package disk

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetAppConfigDir returns the appropriate config directory for the application
func GetAppConfigDir() (string, error) {
	appDir, err := GetAppSupportDir()
	if err != nil {
		return "", fmt.Errorf("failed to get app dir: %w", err)
	}

	configDir := filepath.Join(appDir, ".config")

	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create config dir: %w", err)
	}

	return configDir, nil
}
