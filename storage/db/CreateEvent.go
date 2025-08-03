package db

func (s *Store) CreateEvent(name string, bucket int64) error {
	_, err := s.db.Exec(`INSERT INTO events (name, bucket_id) VALUES (?, ?)`, name, bucket)
	if err != nil {
		return err
	}

	return nil
}
