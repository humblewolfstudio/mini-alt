package storage

import (
	"errors"
	"sync"
	"time"
)

type ChecksumAlgorithm string

const (
	CRC32     ChecksumAlgorithm = "CRC32"
	CRC32C                      = "CRC32C"
	SHA1                        = "SHA1"
	SHA256                      = "SHA256"
	CRC64NVME                   = "CRC64NVME"
)

type ChecksumType string

const (
	COMPOSITE   ChecksumType = "COMPOSITE"
	FULL_OBJECT              = "FULL_OBJECT"
)

type Object struct {
	Id                int64
	ChecksumAlgorithm []ChecksumAlgorithm
	ChecksumType      ChecksumType
	ETag              string
	Key               string
	LastModified      time.Time
	Size              int64
}

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

type Bucket struct {
	Id        int64
	Name      string
	CreatedAt time.Time
}

type InMemoryStore struct {
	mu      sync.Mutex
	buckets map[string]map[string]Object
	meta    map[string]Bucket
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		buckets: make(map[string]map[string]Object),
		meta:    make(map[string]Bucket),
	}
}

func (s *InMemoryStore) PutBucket(bucket string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.buckets[bucket]; exists {
		return errors.New("bucket already exists")
	}
	s.buckets[bucket] = make(map[string]Object)
	s.meta[bucket] = Bucket{Name: bucket, CreatedAt: time.Now()}
	return nil
}

func (s *InMemoryStore) ListBuckets() ([]Bucket, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var buckets []Bucket
	for name, bucket := range s.meta {
		buckets = append(buckets, Bucket{Name: name, CreatedAt: bucket.CreatedAt})
	}
	return buckets, nil
}

func (s *InMemoryStore) PutObject(bucket, objectKey string, size int64) (Object, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	obj := Object{
		Key:          objectKey,
		Size:         size,
		LastModified: time.Now(),
	}

	if s.buckets[bucket] == nil {
		s.buckets[bucket] = make(map[string]Object)
	}

	s.buckets[bucket][objectKey] = obj

	return obj, nil
}

func (s *InMemoryStore) ListObjects(bucket string) ([]Object, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	bucketObjects, exists := s.buckets[bucket]
	if !exists {
		return nil, nil
	}

	objects := make([]Object, 0, len(bucketObjects))
	for _, obj := range bucketObjects {
		objects = append(objects, obj)
	}
	return objects, nil
}

func (s *InMemoryStore) DeleteObject(bucket, objectKey string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if bucketObjects, ok := s.buckets[bucket]; ok {
		delete(bucketObjects, objectKey)
	}

	return nil
}

func (s *InMemoryStore) DeleteBucket(bucketName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.buckets, bucketName)

	return nil
}

func (s *InMemoryStore) GetObject(bucket, key string) (Object, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if bucketObjects, ok := s.buckets[bucket]; ok {
		return bucketObjects[key], nil
	}

	return Object{}, errors.New("the specified key does not exist")
}

func (s *InMemoryStore) GetBucket(bucket string) (Bucket, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if bucket, ok := s.meta[bucket]; ok {
		return bucket, nil
	}

	return Bucket{}, errors.New("the specified bucket does not exist")
}
