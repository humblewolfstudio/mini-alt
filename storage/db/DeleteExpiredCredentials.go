package db

import "log"

func (s *Store) DeleteExpiredCredentials() {
	_, err := s.db.Exec(`DELETE FROM credentials WHERE expires_at IS NOT NULL AND DATE(expires_at) < DATE('now');`)
	if err != nil {
		log.Printf("Failed to delete expired credentials: %v", err)
	} else {
		log.Println("Expired credentials deleted.")
	}
}
