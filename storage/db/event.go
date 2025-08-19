package db

import "mini-alt/models"

func (s *Store) listEvents() ([]models.Event, error) {
	rows, err := s.db.Query(`SElECT * FROM events`)
	if err != nil {
		return nil, err
	}

	var events []models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.Id, &event.Name, &event.Description, &event.BucketId, &event.Endpoint, &event.Token, &event.CreatedAt); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (s *Store) createEvent(name, description, endpoint, token string, bucket int64) error {
	_, err := s.db.Exec(`INSERT INTO events (name, description, endpoint, token, bucket_id) VALUES (?, ?, ?, ?, ?)`, name, description, endpoint, token, bucket)
	if err != nil {
		return err
	}

	return nil
}
