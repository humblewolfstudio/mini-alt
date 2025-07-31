package disk

import "path/filepath"

// GetBucketsDir returns the full path to the data directory
func GetBucketsDir() (string, error) {
	appSupportDir, err := GetAppSupportDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(appSupportDir, "data"), nil
}
