package disk

import (
	"fmt"
	"os"
	"path/filepath"
)

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
