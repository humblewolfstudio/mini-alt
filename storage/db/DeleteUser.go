package db

func (s *Store) DeleteUser(id int64) error {
	user, err := s.GetUserById(id)
	if err != nil {
		return err
	}

	err = s.DeleteCredentials(user.AccessKey)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`DELETE FROM users WHERE id = ?`, id)
	return err
}
