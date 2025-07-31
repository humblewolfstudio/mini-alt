package db

import "mini-alt/models"

func (s *Store) GetUserById(id int64) (models.User, error) {
	row := s.db.QueryRow(`SELECT id, username, password, token, access_key, expires_at FROM users WHERE id = ?`, id)

	var u models.User
	if err := row.Scan(&u.Id, &u.Username, &u.Password, &u.Token, &u.AccessKey, &u.ExpiresAt); err != nil {
		return models.User{}, err
	}

	return u, nil
}
