package disk

import (
	"fmt"
	"os"
	"path/filepath"
)

func DeleteObjectFile(bucketName, objectKey string) error {
	bucketsDir, err := GetBucketsDir()
	if err != nil {
		return err
	}
	path := filepath.Join(bucketsDir, bucketName, objectKey)
	if err := os.RemoveAll(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete object file: %w", err)
	}
	return nil
}
