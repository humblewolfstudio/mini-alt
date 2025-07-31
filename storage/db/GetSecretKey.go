package db

import "mini-alt/utils"

func (s *Store) GetSecretKey(accessKey string) (string, error) {
	row := s.db.QueryRow(`SELECT secret_key_encrypted FROM credentials WHERE access_key = ?`, accessKey)
	var encrypted string
	if err := row.Scan(&encrypted); err != nil {
		return "", err
	}

	return utils.Decrypt(encrypted)
}
