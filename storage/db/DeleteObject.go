package db

func (s *Store) DeleteObject(bucket, key string) error {
	prefix := key + "%"
	_, err := s.db.Exec(`DELETE FROM objects WHERE bucket_name = ? AND key LIKE ? ESCAPE '\'`, bucket, prefix)
	return err
}
