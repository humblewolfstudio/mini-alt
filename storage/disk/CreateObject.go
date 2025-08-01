package disk

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CreateObject(path string, src io.Reader) (int64, error) {
	if info, err := os.Stat(path); err == nil {
		if info.IsDir() {
			return 0, nil
		}
	}

	if src == nil {
		if err := os.MkdirAll(path, 0755); err != nil {
			return 0, fmt.Errorf("failed to create directory: %w", err)
		}

		return 0, nil
	}

	dir := filepath.Dir(path)
	if info, err := os.Stat(dir); err == nil && !info.IsDir() {
		if err := os.Remove(dir); err != nil {
			return 0, fmt.Errorf("cannot convert file to directory: %w", err)
		}
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return 0, fmt.Errorf("failed to create object directory: %w", err)
	}

	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return 0, nil
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
		return 0, nil
	}

	return written, nil
}
