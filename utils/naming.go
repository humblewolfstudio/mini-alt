package utils

import "strings"

func ClearBucketName(bucketName string) string {
	return strings.TrimSuffix(bucketName, "/")
}

func ClearObjectKeyWithBucket(bucket, objectKey string) string {
	return strings.TrimSuffix(bucket, "/") + "/" + strings.TrimPrefix(objectKey, "/")
}
