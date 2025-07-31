package db

import "mini-alt/models"

func (s *Store) ListUsers() ([]models.User, error) {
	rows, err := s.db.Query(`SELECT id, username, expires_at, created_at FROM users`)
	if err != nil {
		return nil, err
	}

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.Id, &u.Username, &u.ExpiresAt, &u.CreatedAt); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}

	return users, nil
}
