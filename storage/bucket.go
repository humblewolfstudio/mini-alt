package storage

import (
	"errors"
	"mini-alt/storage/disk"
)

func (s *Storage) createBucket(bucket string, owner int64) error {
	err := s.store.PutBucket(bucket, owner)
	if err != nil {
		return errors.New("the requested bucket name is not available")
	}

	err = disk.CreateBucket(bucket)
	if err != nil {
		return errors.New("could not create storage directory")
	}

	return nil
}
