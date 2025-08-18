package db

import (
	"time"
)

func (s *Store) PutBucket(bucket string, owner int64) error {
	now := time.Now()
	_, err := s.db.Exec(`INSERT INTO buckets(name, owner, created_at) VALUES (?, ?, ?)`, bucket, owner, now)
	if err != nil {
		return err
	}

	return nil
}
