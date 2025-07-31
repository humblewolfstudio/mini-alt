package db

import "log"

func (s *Store) DeleteExpiredUsers() {
	_, err := s.db.Exec(`DELETE FROM users WHERE expires_at IS NOT NULL AND DATE(expires_at) < DATE('now');`)
	if err != nil {
		log.Printf("Failed to delete expired users: %v", err)
	} else {
		log.Println("Expired users deleted.")
	}
}
