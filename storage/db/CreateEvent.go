package db

func (s *Store) CreateEvent(name, description string, bucket int64) error {
	_, err := s.db.Exec(`INSERT INTO events (name, description, bucket_id) VALUES (?, ?, ?)`, name, description, bucket)
	if err != nil {
		return err
	}

	return nil
}
