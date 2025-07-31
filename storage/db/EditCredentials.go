package db

func (s *Store) EditCredentials(accessKey string, name, description, expiresAt string, status bool) error {
	var expiresAtValue interface{}
	if expiresAt == "" {
		expiresAtValue = nil
	} else {
		expiresAtValue = expiresAt
	}

	_, err := s.db.Exec(`UPDATE credentials SET expires_at = ?, status = ?, name = ?, description = ? WHERE access_key = ?`, expiresAtValue, status, name, description, accessKey)
	return err
}
