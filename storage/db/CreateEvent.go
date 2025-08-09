package db

func (s *Store) CreateEvent(name, description, endpoint, token string, bucket int64, global bool) error {
	_, err := s.db.Exec(`INSERT INTO events (name, description, endpoint, token, bucket_id, global) VALUES (?, ?, ?, ?, ?, ?)`, name, description, endpoint, token, bucket, global)
	if err != nil {
		return err
	}

	return nil
}
