package disk

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func putObject(bucket, key string, r io.Reader) (int64, error) {
	path, err := getSafeObjectPath(bucket, key)
	if err != nil {
		return 0, err
	}

	empty, reader, err := isReaderEmpty(r)
	if err != nil {
		return 0, err
	}

	if empty {
		return 0, os.MkdirAll(path, 0755)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return 0, fmt.Errorf("failed to create object directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return 0, fmt.Errorf("failed to create object file: %w", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	return io.Copy(file, reader)
}

func getObject(bucket, key string) (*os.File, error) {
	path, err := getSafeObjectPath(bucket, key)
	if err != nil {
		return nil, err
	}

	return os.Open(path)
}

func deleteObject(bucket, key string) error {
	path, err := getSafeObjectPath(bucket, key)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(path); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}
