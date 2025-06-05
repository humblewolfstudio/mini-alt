package storage

import (
	"io"
	"os"
	"path/filepath"
)

func CreateBucketDirectory(bucketName string) error {
	path := filepath.Join("uploads" + bucketName)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func DeleteObjectFile(bucketName, objectKey string) {
	path := filepath.Join("uploads", bucketName, objectKey)
	_ = os.Remove(path)
}

func CreateObjectFilePath(bucketName, objectKey string) (string, error) {
	path := filepath.Join("uploads", bucketName, objectKey)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", err
	}

	return path, nil
}

func CreateObject(path string, src io.Reader) (int64, error) {
	file, err := os.Create(path)
	if err != nil {
		return 0, nil
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	written, err := io.Copy(file, src)
	if err != nil {
		return 0, err
	}

	return written, nil
}

func GetObject(bucketName, objectKey string) (string, error) {
	path := filepath.Join("uploads", bucketName, objectKey)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", err
	}

	return path, nil
}
