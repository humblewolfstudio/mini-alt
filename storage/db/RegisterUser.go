package db

import "mini-alt/utils"

func (s *Store) RegisterUser(username, password, accessKey, expiresAt string, admin bool) (int64, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return -1, err
	}
	token := utils.GenerateRandomKey(32)
	hashedToken, err := utils.HashPassword(token)
	if err != nil {
		return -1, err
	}

	var expiresAtValue interface{}
	if expiresAt == "" {
		expiresAtValue = nil
	} else {
		expiresAtValue = expiresAt
	}

	result, err := s.db.Exec(`INSERT INTO users (username, password, token, access_key, admin, expires_at) VALUES (?, ?, ?, ?, ?, ?);`, username, hashedPassword, hashedToken, accessKey, admin, expiresAtValue)
	if err != nil {
		return -1, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return userID, nil
}
