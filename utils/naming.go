package utils

import "strings"

func ClearInput(input string) string {
	return strings.TrimPrefix(input, "/")
}

func ClearObjectKeyWithBucket(bucket, objectKey string) string {
	return bucket + "/" + strings.TrimPrefix(objectKey, "/")
}
