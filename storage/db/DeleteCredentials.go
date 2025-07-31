package db

func (s *Store) DeleteCredentials(accessKey string) error {
	_, err := s.db.Exec(`DELETE FROM credentials WHERE access_key = ?`, accessKey)
	return err
}
