package db

import (
	"database/sql"
	"mini-alt/models"
)

func (s *Store) ListCredentials() ([]models.Credentials, error) {
	rows, err := s.db.Query(`SELECT access_key, expires_at, created_at, status, name, description FROM credentials WHERE user = FALSE`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var credentials []models.Credentials
	for rows.Next() {
		var c models.Credentials
		if err := rows.Scan(&c.AccessKey, &c.ExpiresAt, &c.CreatedAt, &c.Status, &c.Name, &c.Description); err != nil {
			return nil, err
		}

		credentials = append(credentials, c)
	}

	return credentials, nil
}
