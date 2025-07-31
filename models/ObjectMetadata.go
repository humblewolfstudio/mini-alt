package models

import "time"

type ObjectMetadata struct {
	Key                string
	Bucket             string
	CacheControl       string
	ContentDisposition string
	ContentEncoding    string
	ContentLanguage    string
	ContentLength      int64
	ContentMD5         string
	ContentType        string
	Expires            time.Time
}
