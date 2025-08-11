package db

import (
	"mini-alt/models"
)

func (s *Store) GetUserByAccessKey(accessKey string) (models.User, error) {
	row := s.db.QueryRow(`
        SELECT 
            u.id,
            u.username,
            u.password,
            u.token,
            c.access_key,
            u.admin,
            u.expires_at,
            u.created_at
        FROM credentials AS c
        INNER JOIN users AS u ON c.owner = u.id
        WHERE c.access_key = ? AND c.status = 1
    `, accessKey)

	var user models.User
	if err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.Token,
		&user.AccessKey,
		&user.Admin,
		&user.ExpiresAt,
		&user.CreatedAt,
	); err != nil {
		return models.User{}, err
	}

	return user, nil
}
