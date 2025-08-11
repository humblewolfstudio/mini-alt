package db

import (
	"errors"
	"mini-alt/utils"
)

func (s *Store) AuthenticateAdmin(id int64, token string) error {
	user, err := s.GetUserById(id)
	if err != nil {
		return err
	}

	if !user.Admin {
		return errors.New("user is not admin")
	}

	equal := utils.CheckPasswordHash(token, user.Token)

	if !equal {
		return errors.New("invalid token")
	}

	return nil
}
