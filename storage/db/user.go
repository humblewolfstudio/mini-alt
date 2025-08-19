package db

import (
	"errors"
	"log"
	"mini-alt/models"
	"mini-alt/utils"
)

func (s *Store) getUser(username string) (models.User, error) {
	row := s.db.QueryRow(`SELECT id, username, password, token, admin, expires_at FROM users WHERE username = ?`, username)

	var u models.User
	if err := row.Scan(&u.Id, &u.Username, &u.Password, &u.Token, &u.Admin, &u.ExpiresAt); err != nil {
		return models.User{}, err
	}

	return u, nil
}

func (s *Store) getUserById(id int64) (models.User, error) {
	row := s.db.QueryRow(`SELECT id, username, password, token, access_key, admin, expires_at FROM users WHERE id = ?`, id)

	var u models.User
	if err := row.Scan(&u.Id, &u.Username, &u.Password, &u.Token, &u.AccessKey, &u.Admin, &u.ExpiresAt); err != nil {
		return models.User{}, err
	}

	return u, nil
}

func (s *Store) getUserByAccessKey(accessKey string) (models.User, error) {
	row := s.db.QueryRow(`
        SELECT 
            u.id,
            u.username,
            u.password,
            u.token,
            c.access_key,
            u.admin,
            u.expires_at,
            u.created_at
        FROM credentials AS c
        INNER JOIN users AS u ON c.owner = u.id
        WHERE c.access_key = ? AND c.status = 1
    `, accessKey)

	var user models.User
	if err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.Token,
		&user.AccessKey,
		&user.Admin,
		&user.ExpiresAt,
		&user.CreatedAt,
	); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *Store) listUsers() ([]models.User, error) {
	rows, err := s.db.Query(`SELECT id, username, expires_at, created_at FROM users`)
	if err != nil {
		return nil, err
	}

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.Id, &u.Username, &u.ExpiresAt, &u.CreatedAt); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Store) deleteUser(id int64) error {
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

func (s *Store) deleteExpiredUsers() {
	_, err := s.db.Exec(`DELETE FROM users WHERE expires_at IS NOT NULL AND DATE(expires_at) < DATE('now');`)
	if err != nil {
		log.Printf("Failed to delete expired users: %v", err)
	} else {
		log.Println("Expired users deleted.")
	}
}

func (s *Store) authenticateAdmin(id int64, token string) error {
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

func (s *Store) authenticateUser(id int64, token string) error {
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

type LoginResponse struct {
	Id    int64
	Token string
	Admin bool
}

func (s *Store) loginUser(username, password string) (LoginResponse, error) {
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

	return LoginResponse{Id: user.Id, Token: token, Admin: user.Admin}, nil
}

func (s *Store) registerUser(username, password, accessKey, expiresAt string, admin bool) (int64, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return -1, err
	}
	token := utils.GenerateRandomKey(32)
	hashedToken, err := utils.HashPassword(token)
	if err != nil {
		return -1, err
	}

	var expiresAtValue interface{}
	if expiresAt == "" {
		expiresAtValue = nil
	} else {
		expiresAtValue = expiresAt
	}

	result, err := s.db.Exec(`INSERT INTO users (username, password, token, access_key, admin, expires_at) VALUES (?, ?, ?, ?, ?, ?);`, username, hashedPassword, hashedToken, accessKey, admin, expiresAtValue)
	if err != nil {
		return -1, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return userID, nil
}
