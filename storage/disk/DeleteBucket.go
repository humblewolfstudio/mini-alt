package disk

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
