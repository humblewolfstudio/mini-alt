package disk

import "os"

func GetObject(bucketName, objectKey string) (*os.File, error) {
	path, err := GetObjectPath(bucketName, objectKey)
	if err != nil {
		return nil, err
	}
	return os.Open(path)
}
