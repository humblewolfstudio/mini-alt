package db

import (
	"errors"
	"mini-alt/utils"
)

type LoginResponse struct {
	Id    int64
	Token string
}

func (s *Store) LoginUser(username, password string) (LoginResponse, error) {
	user, err := s.GetUser(username)
	if err != nil {
		return LoginResponse{}, err
	}
	equal := utils.CheckPasswordHash(password, user.Password)

	if !equal {
		return LoginResponse{}, errors.New("invalid password")
	}

	token := utils.GenerateRandomKey(32)
	hashedToken, err := utils.HashPassword(token)
	if err != nil {
		return LoginResponse{}, err
	}

	_, err = s.db.Exec(`UPDATE users SET token = ? WHERE username = ?`, hashedToken, username)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{Id: user.Id, Token: token}, nil
}
