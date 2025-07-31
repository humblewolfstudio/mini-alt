package models

import "time"

type Credentials struct {
	Id          int64
	AccessKey   string
	SecretKey   string
	ExpiresAt   *time.Time
	Status      bool
	Name        *string
	Description *string
	CreatedAt   time.Time
}
