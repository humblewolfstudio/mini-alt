package db

func getConfig[T any](s *Store, key string) (T, error) {
	var value T
	err := s.db.QueryRow(`SELECT value FROM config WHERE key = ?`, key).Scan(&value)
	if err != nil {
		var zero T
		return zero, err
	}
	return value, nil
}

func putConfig(s *Store, key string, value interface{}) error {
	_, err := s.db.Exec(`INSERT INTO config (key, value) VALUES (?, ?)`, key, value)
	return err
}
