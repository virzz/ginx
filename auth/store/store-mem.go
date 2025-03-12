package store

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

type MemStore struct {
	keyPrefix string
	maxAge    int

	lock sync.Mutex
	data map[string][]byte
}

func (s *MemStore) WithKeyPrefix(keyPrefix string) Store {
	s.keyPrefix = keyPrefix
	return s
}

func (s *MemStore) WithMaxAge(maxAge int) Store {
	s.maxAge = maxAge
	return s
}

func NewMemStore(opts ...StoreOption) (*MemStore, error) {
	rs := &MemStore{
		keyPrefix: StoreKeyPrefix,
		maxAge:    StoreMaxAge,
		data:      make(map[string][]byte),
	}
	for _, opt := range opts {
		opt(rs)
	}
	return rs, nil
}

func MemStoreWithKeyPrefix(keyPrefix string) func(*MemStore) {
	return func(rs *MemStore) { rs.keyPrefix = keyPrefix }
}

func (s *MemStore) Clear(v Data) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.data, s.keyPrefix+v.Token())
	return nil
}

func (s *MemStore) Get(ctx context.Context, v Data) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if buf, ok := s.data[s.keyPrefix+v.Token()]; ok {
		return json.Unmarshal(buf, v)
	}
	return errors.New("not found")
}

func (s *MemStore) Save(ctx context.Context, v Data, _ ...time.Duration) error {
	if v.Token() == "" {
		v.New()
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	buf, err := json.Marshal(v)
	if err != nil {
		return err
	}
	s.data[s.keyPrefix+v.Token()] = buf
	return nil
}
