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
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		buckets: make(map[string][]Object),
	}
}

func (s *InMemoryStore) CreateBucket(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.buckets[name]; exists {
		return errors.New("bucket already exists")
	}
	s.buckets[name] = make([]Object, 0)
	return nil
}
