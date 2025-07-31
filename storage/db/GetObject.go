package db

import (
	"errors"
	"mini-alt/models"
)

func (s *Store) GetObject(bucket, key string) (models.Object, error) {
	row := s.db.QueryRow(`
		SELECT key, size, last_modified FROM objects WHERE bucket_name = ? AND key = ?`,
		bucket, key)

	var obj models.Object
	obj.Key = key
	if err := row.Scan(&obj.Key, &obj.Size, &obj.LastModified); err != nil {
		return models.Object{}, errors.New("object not found")
	}
	return obj, nil
}
