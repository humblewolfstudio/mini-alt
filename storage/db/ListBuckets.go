package db

import (
	"database/sql"
	"mini-alt/models"
)

func (s *Store) ListBuckets(owner int64) ([]models.Bucket, error) {
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
