package disk

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CreateObjectFilePath(bucketName, objectKey string) (string, error) {
	if strings.Contains(objectKey, "..") || filepath.IsAbs(objectKey) {
		return "", fmt.Errorf("invalid object key: %s", objectKey)
	}

	cleanedKey := filepath.Clean(objectKey)
	bucketsDir, err := GetBucketsDir()
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(bucketsDir, bucketName, cleanedKey)
	bucketRoot := filepath.Join(bucketsDir, bucketName)

	if !strings.HasPrefix(fullPath, bucketRoot) {
		return "", fmt.Errorf("object key escapes bucket path: %s", objectKey)
	}

	dir := filepath.Dir(fullPath)

	if info, err := os.Stat(dir); err == nil && !info.IsDir() {
		if err := os.Remove(dir); err != nil {
			return "", fmt.Errorf("cannot convert file to directory: %w", err)
		}
	}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create object path directories: %w", err)
	}

	return fullPath, nil
}
