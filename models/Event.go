package models

import "time"

type Event struct {
	Id        int64
	Name      string
	BucketId  int64
	CreatedAt time.Time
}
