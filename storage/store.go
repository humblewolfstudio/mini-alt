package storage

import "mini-alt/types"

type Store interface {
	PutObject(bucket, object string, size int64) (Object, error)
	PutBucket(bucket string) error
	PutMetadata(objectId int64, metadata types.Metadata) error
	ListObjects(bucket string) ([]Object, error)
	ListBuckets() ([]Bucket, error)
	GetObject(bucket, key string) (Object, error)
	GetBucket(bucket string) (Bucket, error)
	DeleteObject(bucket, objectKey string) error
	DeleteBucket(bucket string) error
	CreateCredentials(expiresAt string, user bool) (accessKey, secretKey string, err error)
	GetSecretKey(accessKey string) (string, error)
	ListCredentials() ([]Credentials, error)
	DeleteCredentials(accessKey string) error
	DeleteExpiredCredentials()
	RegisterUser(username, password, accessKey, expiresAt string) error
	GetUser(username string) (User, error)
	GetUserById(id int64) (User, error)
	LoginUser(username, password string) (LoginResponse, error)
	AuthenticateUser(id int64, token string) error
	ListUsers() ([]User, error)
	DeleteUser(id int64) error
	DeleteExpiredUsers()
}
