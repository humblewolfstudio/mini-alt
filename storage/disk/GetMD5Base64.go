package disk

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func GetMD5Base64(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash object: %w", err)
	}

	hashSum := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(hashSum), nil
}
