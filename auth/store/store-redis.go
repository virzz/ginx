package store

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/virzz/vlog"
)

type RedisStore struct {
	client    redis.UniversalClient
	keyPrefix string
	maxAge    int
}

func (s *RedisStore) WithKeyPrefix(keyPrefix string) Store {
	s.keyPrefix = keyPrefix
	return s
}

func (s *RedisStore) WithMaxAge(maxAge int) Store {
	s.maxAge = maxAge
	return s
}

func NewRedisStore(client redis.UniversalClient, opts ...StoreOption) (*RedisStore, error) {
	rs := &RedisStore{
		keyPrefix: StoreKeyPrefix,
		maxAge:    StoreMaxAge,
		client:    client,
	}
	for _, opt := range opts {
		opt(rs)
	}
	if rs.client == nil {
		return nil, errors.New("redisstore: client is nil")
	}
	return rs, rs.client.Ping(context.Background()).Err()
}

func (s *RedisStore) Clear(v Data) error {
	return s.client.Del(context.Background(), s.keyPrefix+v.Token()).Err()
}

func (s *RedisStore) Get(ctx context.Context, v Data) error {
	x := s.client.HGetAll(ctx, s.keyPrefix+v.Token())
	if len(x.Val()) == 0 {
		return redis.Nil
	}
	return x.Scan(v)
}

func (s *RedisStore) Save(ctx context.Context, v Data, lifetime ...time.Duration) error {
	if v.Token() == "" || v.Token() == "null" {
		v.New()
	}
	maxAge := time.Duration(s.maxAge) * time.Second
	if len(lifetime) > 0 {
		maxAge = lifetime[0]
	}
	key := s.keyPrefix + v.Token()
	if err := s.client.HSet(ctx, key, v).Err(); err != nil {
		vlog.Error("Failed to hset", "key", key, "err", err.Error())
		return err
	}
	return s.client.Expire(ctx, key, maxAge).Err()
}
