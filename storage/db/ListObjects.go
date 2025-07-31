package db

import (
	"database/sql"
	"mini-alt/models"
)

func (s *Store) ListObjects(bucket string) ([]models.Object, error) {
	rows, err := s.db.Query(`SELECT key, size, last_modified FROM objects WHERE bucket_name = ?`, bucket)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var objs []models.Object
	for rows.Next() {
		var o models.Object
		if err := rows.Scan(&o.Key, &o.Size, &o.LastModified); err != nil {
			return nil, err
		}
		objs = append(objs, o)
	}
	return objs, nil
}
