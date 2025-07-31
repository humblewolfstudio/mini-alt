package db

import "mini-alt/utils"

func (s *Store) CreateCredentials(name, description, expiresAt string, user bool) (accessKey, secretKey string, err error) {
	accessKey = utils.GenerateRandomKey(16)
	secretKey = utils.GenerateRandomKey(32)

	encryptedSecret, err := utils.Encrypt(secretKey)
	if err != nil {
		return "", "", err
	}

	var expiresAtValue interface{}
	if expiresAt == "" {
		expiresAtValue = nil
	} else {
		expiresAtValue = expiresAt
	}

	_, err = s.db.Exec(`INSERT INTO credentials (access_key, secret_key_encrypted, user, expires_at, name, description) VALUES (?, ?, ?, ?, ?, ?)`, accessKey, encryptedSecret, user, expiresAtValue, name, description)
	if err != nil {
		return "", "", err
	}

	return accessKey, secretKey, nil
}
