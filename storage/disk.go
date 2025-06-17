package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// GetAppSupportDir returns the appropriate application support directory for macOS
func GetAppSupportDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	//TODO add support for other systems
	return filepath.Join(home, "Library", "Application Support", "mini-alt"), nil
}

// GetBucketsDir returns the full path to the data directory
func GetBucketsDir() (string, error) {
	appSupportDir, err := GetAppSupportDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(appSupportDir, "data"), nil
}

// EnsureDirectories creates all required directories if they don't exist
func EnsureDirectories() error {
	appSupportDir, err := GetAppSupportDir()
	if err != nil {
		return err
	}

	// Create application support directory if it doesn't exist
	if err := os.MkdirAll(appSupportDir, 0755); err != nil {
		return fmt.Errorf("failed to create application support directory: %w", err)
	}

	// Create buckets directory
	bucketsDir, err := GetBucketsDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(bucketsDir, 0755); err != nil {
		return fmt.Errorf("failed to create buckets directory: %w", err)
	}

	return nil
}

func CreateBucketDirectory(bucketName string) error {
	bucketsDir, err := GetBucketsDir()
	if err != nil {
		return err
	}
	path := filepath.Join(bucketsDir, bucketName)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create bucket directory: %w", err)
	}
	return nil
}

func DeleteObjectFile(bucketName, objectKey string) error {
	bucketsDir, err := GetBucketsDir()
	if err != nil {
		return err
	}
	path := filepath.Join(bucketsDir, bucketName, objectKey)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete object file: %w", err)
	}
	return nil
}

func CreateObjectFilePath(bucketName, objectKey string) (string, error) {
	bucketsDir, err := GetBucketsDir()
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(bucketsDir, bucketName, objectKey)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create object path directories: %w", err)
	}
	return fullPath, nil
}

func CreateObject(path string, src io.Reader) (int64, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return 0, fmt.Errorf("failed to create object directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return 0, fmt.Errorf("failed to create object file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("warning: failed to close file: %v\n", closeErr)
		}
	}()

	written, err := io.Copy(file, src)
	if err != nil {
		return 0, fmt.Errorf("failed to write object content: %w", err)
	}

	return written, nil
}

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

func GetObject(bucketName, objectKey string) (*os.File, error) {
	path, err := GetObjectPath(bucketName, objectKey)
	if err != nil {
		return nil, err
	}
	return os.Open(path)
}

func DeleteBucket(bucketName string) error {
	if bucketName == "" || bucketName == "." || bucketName == "/" {
		return fmt.Errorf("invalid bucket name")
	}

	bucketsDir, err := GetBucketsDir()
	if err != nil {
		return err
	}

	path := filepath.Join(bucketsDir, bucketName)

	absUploads, err := filepath.Abs(bucketsDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for buckets directory: %w", err)
	}

	absTarget, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for target: %w", err)
	}

	if !strings.HasPrefix(absTarget, absUploads) {
		return fmt.Errorf("refusing to delete outside of %s directory", bucketsDir)
	}

	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to delete bucket: %w", err)
	}

	return nil
}
