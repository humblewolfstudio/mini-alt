package models

import "time"

type Event struct {
	Id          int64
	Name        string
	Description string
	BucketId    int64
	Endpoint    string
	Token       string
	CreatedAt   time.Time
}
