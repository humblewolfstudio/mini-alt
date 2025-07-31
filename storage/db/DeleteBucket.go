package db

func (s *Store) DeleteBucket(bucket string) error {
	_, err := s.db.Exec(`DELETE FROM buckets WHERE name = ?`, bucket)
	return err
}
