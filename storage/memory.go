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
	buckets map[string][]Object
	meta    map[string]Bucket
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		buckets: make(map[string][]Object),
		meta:    make(map[string]Bucket),
	}
}

func (s *InMemoryStore) CreateBucket(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.buckets[name]; exists {
		return errors.New("bucket already exists")
	}
	s.buckets[name] = make([]Object, 0)
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

func (s *InMemoryStore) PutObject(bucket, key string, size int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	obj := Object{
		ObjectKey:    key,
		Size:         size,
		LastModified: time.Now(),
	}
	s.buckets[bucket] = append(s.buckets[bucket], obj)
}
