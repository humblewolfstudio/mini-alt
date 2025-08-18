package db

import (
	"database/sql"
	"log"
	"mini-alt/models"
	"mini-alt/utils"
)

func (s *Store) getSecretKey(accessKey string) (string, error) {
	row := s.db.QueryRow(`SELECT secret_key_encrypted FROM credentials WHERE access_key = ?`, accessKey)
	var encrypted string
	if err := row.Scan(&encrypted); err != nil {
		return "", err
	}

	return utils.Decrypt(encrypted)
}

func (s *Store) listCredentials(owner int64) ([]models.Credentials, error) {
	rows, err := s.db.Query(`SELECT access_key, expires_at, created_at, status, name, description FROM credentials WHERE user = FALSE AND owner = ?`, owner)
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

func (s *Store) putCredentials(name, description, expiresAt string, user bool, owner int64) (accessKey, secretKey string, err error) {
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

	_, err = s.db.Exec(`INSERT INTO credentials (access_key, secret_key_encrypted, user, expires_at, name, description, owner) VALUES (?, ?, ?, ?, ?, ?, ?)`, accessKey, encryptedSecret, user, expiresAtValue, name, description, owner)
	if err != nil {
		return "", "", err
	}

	return accessKey, secretKey, nil
}

func (s *Store) deleteCredentials(accessKey string) error {
	_, err := s.db.Exec(`DELETE FROM credentials WHERE access_key = ?`, accessKey)
	return err
}

func (s *Store) deleteExpiredCredentials() {
	_, err := s.db.Exec(`DELETE FROM credentials WHERE expires_at IS NOT NULL AND DATE(expires_at) < DATE('now');`)
	if err != nil {
		log.Printf("Failed to delete expired credentials: %v", err)
	} else {
		log.Println("Expired credentials deleted.")
	}
}

func (s *Store) addCredentialsOwner(accessKey string, owner int64) error {
	_, err := s.db.Exec(`UPDATE credentials SET owner = ? WHERE access_key = ?;`, owner, accessKey)
	return err
}

func (s *Store) editCredentials(accessKey string, name, description, expiresAt string, status bool) error {
	var expiresAtValue interface{}
	if expiresAt == "" {
		expiresAtValue = nil
	} else {
		expiresAtValue = expiresAt
	}

	_, err := s.db.Exec(`UPDATE credentials SET expires_at = ?, status = ?, name = ?, description = ? WHERE access_key = ?`, expiresAtValue, status, name, description, accessKey)
	return err
}
