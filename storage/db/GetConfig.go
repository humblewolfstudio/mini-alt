package db

func GetConfig[T any](s *Store, key string) (T, error) {
	var value T
	err := s.db.QueryRow(`SELECT value FROM config WHERE key = ?`, key).Scan(&value)
	if err != nil {
		var zero T
		return zero, err
	}
	return value, nil
}
