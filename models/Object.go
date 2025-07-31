package models

import "time"

type Object struct {
	Id           int64
	ETag         string
	Key          string
	LastModified time.Time
	Size         int64
}
