package disk

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CreateObject(path string, src io.Reader) (int64, error) {
	empty, reader, err := isReaderEmpty(src)
	if err != nil {
		return 0, err
	}

	if empty {
		if err := os.MkdirAll(path, 0755); err != nil {
			return 0, fmt.Errorf("failed to create directory: %w", err)
		}

		return 0, nil
	}

	if info, err := os.Stat(path); err == nil {
		if info.IsDir() {
			return 0, nil
		}
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

	written, err := io.Copy(file, reader)
	if err != nil {
		return 0, nil
	}

	return written, nil
}

func isReaderEmpty(r io.Reader) (bool, io.Reader, error) {
	if f, ok := r.(*os.File); ok {
		info, err := f.Stat()
		if err != nil {
			return false, r, err
		}
		if info.IsDir() {
			return true, r, nil
		}
	}

	br := bufio.NewReader(r)
	_, err := br.Peek(1)
	if err != nil {
		if err == io.EOF {
			return true, br, nil
		}
		return false, br, err
	}
	return false, br, nil
}
