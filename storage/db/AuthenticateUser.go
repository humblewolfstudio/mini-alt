package db

import (
	"errors"
	"mini-alt/utils"
)

func (s *Store) AuthenticateUser(id int64, token string) error {
	user, err := s.GetUserById(id)
	if err != nil {
		return err
	}

	equal := utils.CheckPasswordHash(token, user.Token)

	if !equal {
		return errors.New("invalid token")
	}

	return nil
}
