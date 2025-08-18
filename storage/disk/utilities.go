package disk

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func isReaderEmpty(r io.Reader) (bool, io.Reader, error) {
	if f, ok := r.(*os.File); ok {
		info, err := f.Stat()
		if err != nil {
			return false, r, err
		}
		if info.IsDir() {
			return true, r, nil
		}
	}

	br := bufio.NewReader(r)
	_, err := br.Peek(1)
	if err != nil {
		if err == io.EOF {
			return true, br, nil
		}
		return false, br, err
	}
	return false, br, nil
}

func getSafeObjectPath(bucket, key string) (string, error) {
	if key == "" || key == "." {
		return "", fmt.Errorf("invalid empty object key")
	}

	if filepath.IsAbs(key) {
		return "", fmt.Errorf("invalid absolute object key: %s", key)
	}

	cleanedKey := filepath.Clean(key)

	rootDir, err := getBucketsPath()
	if err != nil {
		return "", err
	}

	bucketRoot := filepath.Join(rootDir, bucket)
	fullPath := filepath.Join(bucketRoot, cleanedKey)

	rel, err := filepath.Rel(bucketRoot, fullPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("object key escapes bucket path: %s", key)
	}

	return fullPath, nil
}

// getAppConfigPath returns the appropriate config directory for the application
func getAppConfigPath() (string, error) {
	appDir, err := getAppSupportPath()
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

// getBucketsPath returns the full path to the data directory
func getBucketsPath() (string, error) {
	appSupportDir, err := getAppSupportPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(appSupportDir, "data"), nil
}

// getAppSupportPath returns the appropriate application support directory based on the OS.
func getAppSupportPath() (string, error) {
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
