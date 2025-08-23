package storage

import (
	"io"
	"mini-alt/models"
	"mini-alt/storage/db"
	"mini-alt/types"
	"mini-alt/utils"
)

type Storage struct {
	store *db.Store
}

// ----------------- Management --------------

func NewStorage(store *db.Store) *Storage {
	return &Storage{store: store}
}

// ----------------- Buckets -----------------

func (s *Storage) PutBucket(bucket string, owner int64) (bool, utils.Error) {
	return s.putBucket(bucket, owner)
}

func (s *Storage) DeleteBucket(bucket string) (bool, utils.Error) { return s.deleteBucket(bucket) }

func (s *Storage) HeadBucket(bucket string) (bool, utils.Error) { return s.headBucket(bucket) }

func (s *Storage) ListBuckets(owner int64) []models.Bucket {
	return s.listBuckets(owner)
}

// ----------------- Objects -----------------

func (s *Storage) PutObject(bucket, objectKey string, body io.Reader, metadata types.Metadata, owner int64) (string, utils.Error) {
	return s.putObject(bucket, objectKey, body, metadata, owner)
}

func (s *Storage) DeleteObject(bucket, objectKey string) (bool, utils.Error) {
	return s.deleteObject(bucket, objectKey)
}

func (s *Storage) HeadObject(bucket, objectKey string) (*models.Object, *models.ObjectMetadata, utils.Error) {
	return s.headObject(bucket, objectKey)
}

func (s *Storage) CopyObject(srcBucket, srcKey, dstBucket, dstKey string) (*models.Object, utils.Error) {
	return s.copyObject(srcBucket, srcKey, dstBucket, dstKey)
}

func (s *Storage) GetObjectPath(bucket, objectKey string) (string, utils.Error) {
	return s.getObjectPath(bucket, objectKey)
}
