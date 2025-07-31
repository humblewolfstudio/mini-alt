package db

import "mini-alt/utils"

func (s *Store) RegisterUser(username, password, accessKey, expiresAt string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	token := utils.GenerateRandomKey(32)
	hashedToken, err := utils.HashPassword(token)
	if err != nil {
		return err
	}

	var expiresAtValue interface{}
	if expiresAt == "" {
		expiresAtValue = nil
	} else {
		expiresAtValue = expiresAt
	}

	_, err = s.db.Exec(`INSERT INTO users (username, password, token, access_key, expires_at) VALUES (?, ?, ?, ?, ?);`, username, hashedPassword, hashedToken, accessKey, expiresAtValue)
	if err != nil {
		return err
	}

	return nil
}
