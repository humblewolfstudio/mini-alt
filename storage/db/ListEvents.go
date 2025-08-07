package db

import (
	"mini-alt/models"
)

func (s *Store) ListEvents() ([]models.Event, error) {
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
