package db

import (
	"mini-alt/models"
	"time"
)

func (s *Store) PutObject(bucket, object string, size int64) (models.Object, error) {
	now := time.Now()
	row := s.db.QueryRow(`
	INSERT INTO objects(bucket_name, key, size, last_modified)
	VALUES (?, ?, ?, ?)
	ON CONFLICT(bucket_name, key) DO UPDATE 
	SET size = excluded.size, last_modified = excluded.last_modified
	RETURNING id`,
		bucket, object, size, now)

	var id int64
	if err := row.Scan(&id); err != nil {
		return models.Object{}, err
	}

	return models.Object{Id: id, Key: object, Size: size, LastModified: now}, nil
}
