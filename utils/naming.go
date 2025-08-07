package utils

import "strings"

func ClearInput(input string) string {
	return strings.TrimPrefix(input, "/")
}

func ClearObjectKeyWithBucket(bucket, objectKey string) string {
	return strings.TrimSuffix(bucket, "/") + "/" + strings.TrimPrefix(objectKey, "/")
}
