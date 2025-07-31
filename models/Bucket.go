package models

import "time"

type Bucket struct {
	Id            int64
	Name          string
	NumberObjects int64
	Size          int64
	CreatedAt     time.Time
}
