package db

import "mini-alt/models"

func (s *Store) GetBucketById(bucketId int64) (models.Bucket, error) {
	row := s.db.QueryRow(`SELECT * FROM buckets WHERE id = ?`, bucketId)

	var b models.Bucket
	if err := row.Scan(&b.Id, &b.Name, &b.CreatedAt); err != nil {
		return models.Bucket{}, err
	}

	return b, nil
}
