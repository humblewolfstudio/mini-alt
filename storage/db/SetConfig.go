package db

func SetConfig(s *Store, key string, value interface{}) error {
	_, err := s.db.Exec(`INSERT INTO config (key, value) VALUES (?, ?)`, key, value)
	return err
}
