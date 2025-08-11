package db

func (s *Store) CreateEvent(name, description, endpoint, token string, bucket int64) error {
	_, err := s.db.Exec(`INSERT INTO events (name, description, endpoint, token, bucket_id) VALUES (?, ?, ?, ?, ?)`, name, description, endpoint, token, bucket)
	if err != nil {
		return err
	}

	return nil
}
