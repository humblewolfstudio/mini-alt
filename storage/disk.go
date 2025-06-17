package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const BucketsDir = "data"

func CreateBucketDirectory(bucketName string) error {
	path := filepath.Join(BucketsDir, bucketName)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func DeleteObjectFile(bucketName, objectKey string) {
	path := filepath.Join(BucketsDir, bucketName, objectKey)
	_ = os.Remove(path)
}

func CreateObjectFilePath(bucketName, objectKey string) (string, error) {
	fullPath := filepath.Join(BucketsDir, bucketName, objectKey)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}

	return fullPath, nil
}

func CreateObject(path string, src io.Reader) (int64, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return 0, err
	}

	file, err := os.Create(path)
	if err != nil {
		return 0, err
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

func GetObjectPath(bucketName, objectKey string) (string, error) {
	path := filepath.Join(BucketsDir, bucketName, objectKey)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path, err
	}

	return path, nil
}

func GetObject(bucketName, objectKey string) (io.ReadWriter, error) {
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

	path := filepath.Join(BucketsDir, bucketName)

	absUploads, _ := filepath.Abs(BucketsDir)
	absTarget, _ := filepath.Abs(path)
	if !strings.HasPrefix(absTarget, absUploads) {
		return fmt.Errorf("refusing to delete outside of %s directory", BucketsDir)
	}

	if err := os.RemoveAll(path); err != nil {
		return err
	}

	return nil
}
