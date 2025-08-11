package db

func (s *Store) AddBucketOwner(bucketName string, bucketOwner int64) error {
	_, err := s.db.Exec(`UPDATE buckets SET owner = ? WHERE name = ?;`, bucketOwner, bucketName)
	return err
}
