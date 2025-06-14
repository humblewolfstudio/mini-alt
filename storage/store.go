package storage

import "mini-alt/types"

type Store interface {
	PutObject(bucket, object string, size int64) (Object, error)
	PutBucket(bucket string) error
	PutMetadata(objectId int64, metadata types.Metadata) error
	ListObjects(bucket string) ([]Object, error)
	ListBuckets() ([]Bucket, error)
	GetObject(bucket, key string) (Object, error)
	GetBucket(bucket string) (Bucket, error)
	DeleteObject(bucket, objectKey string) error
	DeleteBucket(bucket string) error
}
