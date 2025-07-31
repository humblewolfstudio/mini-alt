package disk

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetObjectPath(bucketName, objectKey string) (string, error) {
	bucketsDir, err := GetBucketsDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(bucketsDir, bucketName, objectKey)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path, fmt.Errorf("object not found: %w", err)
	}
	return path, nil
}
