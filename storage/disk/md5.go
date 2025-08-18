package disk

import (
	"crypto/md5"
	"encoding/base64"
	"io"
	"os"
	"path/filepath"
)

func getMD5Base64(bucket, key string) (string, error) {
	rootDir, err := getBucketsPath()
	if err != nil {
		return "", err
	}

	path := filepath.Join(rootDir, bucket, key)

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashSum := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(hashSum), nil
}
