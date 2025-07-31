package disk

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// GetAppSupportDir returns the appropriate application support directory based on the OS.
func GetAppSupportDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	var baseDir string
	appName := "mini-alt"

	switch runtime.GOOS {
	case "darwin":
		baseDir = filepath.Join(home, "Library", "Application Support", appName)

	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("APPDATA environment variable is not set")
		}
		baseDir = filepath.Join(appData, appName)

	case "linux":
		xdgConfig := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfig == "" {
			xdgConfig = filepath.Join(home, ".config")
		}
		baseDir = filepath.Join(xdgConfig, appName)

	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	return baseDir, nil
}
