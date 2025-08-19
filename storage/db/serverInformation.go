package db

type ServerInformation struct {
	NumberBuckets int64
	NumberObjects int64
	Usage         int64
}

func (s *Store) getServerInformation() (*ServerInformation, error) {
	query := `
		SELECT
			COUNT(DISTINCT b.id) AS total_buckets,
			COUNT(o.id) AS total_objects,
			IFNULL(SUM(o.size), 0) AS total_size
		FROM 
			buckets b
		LEFT JOIN 
			objects o ON o.bucket_name = b.name;
	`

	row := s.db.QueryRow(query)

	var server ServerInformation
	if err := row.Scan(&server.NumberBuckets, &server.NumberObjects, &server.Usage); err != nil {
		return nil, err
	}

	return &server, nil
}
