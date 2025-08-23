package storage

import (
	"mini-alt/models"
	"mini-alt/storage/disk"
	"mini-alt/utils"
)

func (s *Storage) putBucket(bucket string, owner int64) (bool, utils.Error) {
	err := s.store.PutBucket(bucket, owner)
	if err != nil {
		return false, utils.BucketAlreadyExists
	}

	err = disk.CreateBucket(bucket)
	if err != nil {
		return false, utils.FailedToCreateBucket
	}

	return true, ""
}

func (s *Storage) deleteBucket(bucket string) (bool, utils.Error) {
	_, err := s.store.GetBucket(bucket)
	if err != nil {
		return false, utils.NoSuchBucket
	}

	hasObjects, err := s.store.BucketHasObjects(bucket)
	if err != nil {
		return false, utils.InvalidRequest
	}

	if hasObjects {
		return false, utils.BucketIsNotEmpty
	}

	err = disk.DeleteBucket(bucket)
	if err != nil {
		return false, utils.FailedToDeleteBucketDirectory
	}

	err = s.store.DeleteBucket(bucket)
	if err != nil {
		return false, utils.FailedToDeleteBucket
	}

	return true, ""
}

func (s *Storage) headBucket(bucket string) (bool, utils.Error) {
	_, err := s.store.GetBucket(bucket)
	if err != nil {
		return false, utils.NoSuchBucket
	}

	return true, ""
}

func (s *Storage) listBuckets(owner int64) []models.Bucket {
	buckets, _ := s.store.ListBuckets(owner)

	return buckets
}
