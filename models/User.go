package models

import "time"

type User struct {
	Id        int64
	Username  string
	Password  string
	Token     string
	AccessKey string
	ExpiresAt *time.Time
	CreatedAt time.Time
}
