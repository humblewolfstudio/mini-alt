package storage

import (
	"mini-alt/models"
	"mini-alt/types"
)

type Store interface {
	PutObject(bucket, object string, size int64) (models.Object, error)
	PutBucket(bucket string) error
	PutMetadata(objectId int64, metadata types.Metadata) error
	ListObjects(bucket string) ([]models.Object, error)
	ListBuckets() ([]models.Bucket, error)
	GetObject(bucket, key string) (models.Object, error)
	GetBucket(bucket string) (models.Bucket, error)
	DeleteObject(bucket, objectKey string) error
	DeleteBucket(bucket string) error
	CreateCredentials(name, description, expiresAt string, user bool) (accessKey, secretKey string, err error)
	EditCredentials(accessKey string, name, description, expiresAt string, status bool) error
	GetSecretKey(accessKey string) (string, error)
	ListCredentials() ([]models.Credentials, error)
	DeleteCredentials(accessKey string) error
	DeleteExpiredCredentials()
	RegisterUser(username, password, accessKey, expiresAt string) error
	GetUser(username string) (models.User, error)
	GetUserById(id int64) (models.User, error)
	LoginUser(username, password string) (LoginResponse, error)
	AuthenticateUser(id int64, token string) error
	ListUsers() ([]models.User, error)
	DeleteUser(id int64) error
	DeleteExpiredUsers()
}
