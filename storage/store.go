package storage

type Store interface {
	ListBuckets() []Bucket
	CreateBucket(name string) error
	PutObject(bucket, object string, size int64)
}
