package db

import (
	"database/sql"
	"mini-alt/models"
	"mini-alt/types"
)

type Store struct {
	db *sql.DB
}

// ----------------- Management --------------

func NewStore(path string) (*Store, error) {
	return newSQLiteStore(path)
}

// ----------------- Buckets -----------------

func (s *Store) GetBucket(bucket string) (models.Bucket, error) {
	return s.getBucket(bucket)
}

func (s *Store) GetBucketById(id int64) (models.Bucket, error) {
	return s.getBucketById(id)
}

func (s *Store) PutBucket(bucket string, owner int64) error {
	return s.putBucket(bucket, owner)
}

func (s *Store) ListBuckets(owner int64) ([]models.Bucket, error) {
	return s.listBuckets(owner)
}

func (s *Store) DeleteBucket(bucket string) error {
	return s.deleteBucket(bucket)
}

func (s *Store) AddBucketOwner(bucketName string, bucketOwner int64) error {
	return s.addBucketOwner(bucketName, bucketOwner)
}

// ----------------- Users -------------------

func (s *Store) GetUser(username string) (models.User, error) {
	return s.getUser(username)
}

func (s *Store) GetUserById(id int64) (models.User, error) {
	return s.getUserById(id)
}

func (s *Store) GetUserByAccessKey(accessKey string) (models.User, error) {
	return s.getUserByAccessKey(accessKey)
}

func (s *Store) ListUsers() ([]models.User, error) {
	return s.listUsers()
}

func (s *Store) DeleteUser(id int64) error {
	return s.deleteUser(id)
}

func (s *Store) DeleteExpiredUsers() {
	s.deleteExpiredUsers()
}

func (s *Store) AuthenticateAdmin(id int64, token string) error {
	return s.authenticateAdmin(id, token)
}

func (s *Store) AuthenticateUser(id int64, token string) error {
	return s.authenticateUser(id, token)
}

func (s *Store) LoginUser(username, password string) (LoginResponse, error) {
	return s.loginUser(username, password)
}

func (s *Store) RegisterUser(username, password, accessKey, expiresAt string, admin bool) (int64, error) {
	return s.registerUser(username, password, accessKey, expiresAt, admin)
}

// ----------------- Objects -----------------

func (s *Store) GetObject(bucket, key string) (models.Object, error) {
	return s.getObject(bucket, key)
}

func (s *Store) ListObjects(bucket string) ([]models.Object, error) {
	return s.listObjects(bucket)
}

func (s *Store) PutObject(bucket, object string, size int64) (models.Object, error) {
	return s.putObject(bucket, object, size)
}

func (s *Store) DeleteObject(bucket, key string) error {
	return s.deleteObject(bucket, key)
}

// ----------------- Metadata ----------------

func (s *Store) GetMetadata(objectId int64) (models.ObjectMetadata, error) {
	return s.getMetadata(objectId)
}

func (s *Store) PutMetadata(objectId int64, metadata types.Metadata) error {
	return s.putMetadata(objectId, metadata)
}

func (s *Store) CopyMetadata(oldObjectId, newObjectId int64) error {
	return s.copyMetadata(oldObjectId, newObjectId)
}

// ----------------- Credentials -------------

func (s *Store) GetSecretKey(accessKey string) (string, error) {
	return s.getSecretKey(accessKey)
}

func (s *Store) ListCredentials(owner int64) ([]models.Credentials, error) {
	return s.listCredentials(owner)
}

func (s *Store) PutCredentials(name, description, expiresAt string, user bool, owner int64) (accessKey, secretKey string, err error) {
	return s.putCredentials(name, description, expiresAt, user, owner)
}

func (s *Store) DeleteCredentials(accessKey string) error {
	return s.deleteCredentials(accessKey)
}

func (s *Store) DeleteExpiredCredentials() {
	s.deleteExpiredCredentials()
}

func (s *Store) AddCredentialsOwner(accessKey string, owner int64) error {
	return s.addCredentialsOwner(accessKey, owner)
}

func (s *Store) EditCredentials(accessKey string, name, description, expiresAt string, status bool) error {
	return s.editCredentials(accessKey, name, description, expiresAt, status)
}

// ----------------- Events ------------------

func (s *Store) ListEvents() ([]models.Event, error) {
	return s.listEvents()
}

func (s *Store) CreateEvent(name, description, endpoint, token string, bucket int64) error {
	return s.createEvent(name, description, endpoint, token, bucket)
}

// ----------------- Config ------------------

func GetConfig[T any](s *Store, key string) (T, error) {
	return getConfig[T](s, key)
}

func PutConfig(s *Store, key string, value interface{}) error {
	return putConfig(s, key, value)
}

// ----------------- Server Information ------

func (s *Store) GetServerInformation() (*ServerInformation, error) {
	return s.getServerInformation()
}
