package db

import (
	"database/sql"
	"errors"
	"mini-alt/models"
	"time"
)

func (s *Store) getObject(bucket, key string) (models.Object, error) {
	row := s.db.QueryRow(`
		SELECT id, key, size, last_modified FROM objects WHERE bucket_name = ? AND key = ?`,
		bucket, key)

	var obj models.Object
	obj.Key = key
	if err := row.Scan(&obj.Id, &obj.Key, &obj.Size, &obj.LastModified); err != nil {
		return models.Object{}, errors.New("object not found")
	}
	return obj, nil
}

func (s *Store) listObjects(bucket string) ([]models.Object, error) {
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

func (s *Store) putObject(bucket, object string, size int64) (models.Object, error) {
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

func (s *Store) deleteObject(bucket, key string) error {
	prefix := key + "%"
	_, err := s.db.Exec(`DELETE FROM objects WHERE bucket_name = ? AND key LIKE ? ESCAPE '\'`, bucket, prefix)
	return err
}
