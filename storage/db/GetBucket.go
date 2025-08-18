package db

import "mini-alt/models"

func (s *Store) GetBucket(bucket string) (models.Bucket, error) {
	row := s.db.QueryRow(`SELECT * FROM buckets WHERE name = ?`, bucket)
	var b models.Bucket
	b.Name = bucket
	if err := row.Scan(&b.Id, &b.Name, &b.Owner, &b.CreatedAt); err != nil {
		return models.Bucket{}, err
	}

	return b, nil
}
