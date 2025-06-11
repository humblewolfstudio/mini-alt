package storage

import (
	"errors"
	"sync"
	"time"
)

type Object struct {
	Id              int64
	BucketName      string
	ObjectKey       string
	Size            int64
	ContentType     string
	LastModified    time.Time
	IsDeletedMarker bool
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

func (s *InMemoryStore) CreateBucket(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.buckets[name]; exists {
		return errors.New("bucket already exists")
	}
	s.buckets[name] = make(map[string]Object)
	s.meta[name] = Bucket{Name: name, CreatedAt: time.Now()}
	return nil
}

func (s *InMemoryStore) ListBuckets() []Bucket {
	s.mu.Lock()
	defer s.mu.Unlock()

	var buckets []Bucket
	for name, bucket := range s.meta {
		buckets = append(buckets, Bucket{Name: name, CreatedAt: bucket.CreatedAt})
	}
	return buckets
}

func (s *InMemoryStore) PutObject(bucket, objectKey string, size int64) Object {
	s.mu.Lock()
	defer s.mu.Unlock()

	obj := Object{
		ObjectKey:    objectKey,
		Size:         size,
		LastModified: time.Now(),
	}

	if s.buckets[bucket] == nil {
		s.buckets[bucket] = make(map[string]Object)
	}

	s.buckets[bucket][objectKey] = obj

	return obj
}

func (s *InMemoryStore) ListObjects(bucket string) []Object {
	s.mu.Lock()
	defer s.mu.Unlock()

	bucketObjects, exists := s.buckets[bucket]
	if !exists {
		return nil
	}

	objects := make([]Object, 0, len(bucketObjects))
	for _, obj := range bucketObjects {
		objects = append(objects, obj)
	}
	return objects
}

func (s *InMemoryStore) DeleteObject(bucket, objectKey string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if bucketObjects, ok := s.buckets[bucket]; ok {
		delete(bucketObjects, objectKey)
	}
}

func (s *InMemoryStore) DeleteBucket(bucketName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.buckets, bucketName)
}
