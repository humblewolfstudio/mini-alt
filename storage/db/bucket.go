package db

import (
	"database/sql"
	"mini-alt/models"
)

func (s *Store) getBucket(bucket string) (models.Bucket, error) {
	row := s.db.QueryRow(`SELECT * FROM buckets WHERE name = ?`, bucket)
	var b models.Bucket
	b.Name = bucket
	if err := row.Scan(&b.Id, &b.Name, &b.Owner, &b.CreatedAt); err != nil {
		return models.Bucket{}, err
	}

	return b, nil
}

func (s *Store) getBucketById(id int64) (models.Bucket, error) {
	row := s.db.QueryRow(`SELECT * FROM buckets WHERE id = ?`, id)

	var b models.Bucket
	if err := row.Scan(&b.Id, &b.Name, &b.Owner, &b.CreatedAt); err != nil {
		return models.Bucket{}, err
	}

	return b, nil
}

func (s *Store) putBucket(bucket string, owner int64) error {
	_, err := s.db.Exec(`INSERT INTO buckets(name, owner) VALUES (?, ?)`, bucket, owner)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) listBuckets(owner int64) ([]models.Bucket, error) {
	query := `
		SELECT 
			b.id,
			b.name,
			b.created_at,
			IFNULL(COUNT(o.id), 0) AS number_objects,
			IFNULL(SUM(o.size), 0) AS total_size
		FROM 
			buckets b
		LEFT JOIN 
			objects o ON o.bucket_name = b.name
		WHERE 
			b.owner = ? 
			OR EXISTS (
				SELECT 1 
				FROM users u 
				WHERE u.id = ? AND u.admin = 1
			)
		GROUP BY 
			b.id, b.name, b.created_at
	`

	rows, err := s.db.Query(query, owner, owner)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var buckets []models.Bucket
	for rows.Next() {
		var b models.Bucket
		if err := rows.Scan(&b.Id, &b.Name, &b.CreatedAt, &b.NumberObjects, &b.Size); err != nil {
			return nil, err
		}
		buckets = append(buckets, b)
	}
	return buckets, nil
}

func (s *Store) deleteBucket(bucket string) error {
	_, err := s.db.Exec(`DELETE FROM buckets WHERE name = ?`, bucket)
	return err
}

func (s *Store) addBucketOwner(bucketName string, bucketOwner int64) error {
	_, err := s.db.Exec(`UPDATE buckets SET owner = ? WHERE name = ?;`, bucketOwner, bucketName)
	return err
}

func (s *Store) bucketHasObjects(bucket string) (bool, error) {
	var exists bool
	err := s.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM objects WHERE bucket_name = ?)`, bucket).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
