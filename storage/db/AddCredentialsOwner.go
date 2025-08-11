package db

func (s *Store) AddCredentialsOwner(accessKey string, owner int64) error {
	_, err := s.db.Exec(`UPDATE credentials SET owner = ? WHERE access_key = ?;`, owner, accessKey)
	return err
}
