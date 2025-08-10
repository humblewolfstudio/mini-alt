package db

import "mini-alt/models"

func (s *Store) GetUser(username string) (models.User, error) {
	row := s.db.QueryRow(`SELECT id, username, password, token, admin, expires_at FROM users WHERE username = ?`, username)

	var u models.User
	if err := row.Scan(&u.Id, &u.Username, &u.Password, &u.Token, &u.Admin, &u.ExpiresAt); err != nil {
		return models.User{}, err
	}

	return u, nil
}
