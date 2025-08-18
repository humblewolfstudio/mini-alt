package disk

import (
	"fmt"
	"io"
	"os"
)

// ----------------- Buckets -----------------

func CreateBucket(bucket string) error {
	return createBucket(bucket)
}

func DeleteBucket(bucket string) error {
	// TODO: confirm the name in the API endpoint...
	if bucket == "" || bucket == "." || bucket == "/" {
		return fmt.Errorf("invalid bucket name")
	}

	return deleteBucket(bucket)
}

// ----------------- Objects -----------------

func PutObject(bucket string, key string, r io.Reader) (int64, error) {
	return putObject(bucket, key, r)
}

func GetObject(bucket string, key string) (*os.File, error) {
	return getObject(bucket, key)
}

func DeleteObject(bucket string, key string) error {
	return deleteObject(bucket, key)
}

// ----------------- MD5 ---------------------

func GetMD5Base64(bucket string, key string) (string, error) {
	return getMD5Base64(bucket, key)
}

// ----------------- System ------------------

func GetSystemSpecs() (SystemSpecs, error) {
	return getSystemSpecs()
}

// ----------------- Path --------------------

func GetSafeObjectPath(bucket string, key string) (string, error) {
	return getSafeObjectPath(bucket, key)
}

func GetAppConfigPath() (string, error) {
	return getAppConfigPath()
}

func GetBucketsPath() (string, error) {
	return getBucketsPath()
}
