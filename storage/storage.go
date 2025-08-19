package storage

import (
	"mini-alt/storage/db"
)

type Storage struct {
	store *db.Store
}

// ----------------- Management --------------

func NewStorage(store *db.Store) *Storage {
	return &Storage{store: store}
}

// ----------------- Buckets -----------------

func (s *Storage) CreateBucket(bucketName string, owner int64) error {
	return s.createBucket(bucketName, owner)
}
