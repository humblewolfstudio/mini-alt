package storage

type Store interface {
	ListBuckets() []Bucket
	CreateBucket(name string) error
	PutObject(bucket, object string, size int64) Object
	ListObjects(bucket string) []Object
	DeleteObject(bucket, objectKey string)
	DeleteBucket(bucket string)
	GetObject(bucket, object string) (Object, error)
}
