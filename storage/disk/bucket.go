package disk

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func createBucketDir(bucket string) error {
	rootDir, err := getBucketsPath()
	if err != nil {
		return err
	}

	path := filepath.Join(rootDir, bucket)
	err = os.MkdirAll(path, os.ModePerm)

	return err
}

func deleteBucketDir(bucket string) error {
	rootDir, err := getBucketsPath()
	if err != nil {
		return err
	}

	deletePath := filepath.Join(rootDir, bucket)
	absoluteRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for buckets directory: %w", err)
	}

	absoluteDeleteDir, err := filepath.Abs(deletePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for target: %w", err)
	}

	if !strings.HasPrefix(absoluteDeleteDir, absoluteRootDir) {
		return fmt.Errorf("refusing to delete outside of %s directory", rootDir)
	}

	err = os.RemoveAll(deletePath)

	return err
}
