package db

import (
	"time"
)

func (s *Store) PutBucket(bucket string) error {
	now := time.Now()
	_, err := s.db.Exec(`INSERT INTO buckets(name, created_at) VALUES (?, ?)`, bucket, now)
	if err != nil {
		return err
	}

	return nil
}
