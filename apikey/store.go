package apikey

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store interface {
	Get(context.Context, Data) error
	Save(context.Context, Data, ...time.Duration) error
	Clear(Data) error
}

type RedisStore struct {
	client    redis.UniversalClient // client to connect to redis
	keyPrefix string                // key prefix with which the session will be stored
	maxAge    int
}

type Option func(*RedisStore)

func WithClient(client redis.UniversalClient) func(*RedisStore) {
	return func(rs *RedisStore) { rs.client = client }
}

func WithKeyPrefix(keyPrefix string) func(*RedisStore) {
	return func(rs *RedisStore) { rs.keyPrefix = keyPrefix }
}

func WithMaxAge(maxAge int) func(*RedisStore) {
	return func(rs *RedisStore) { rs.maxAge = maxAge }
}

func NewRedisStore(opts ...Option) (*RedisStore, error) {
	rs := &RedisStore{
		keyPrefix: "ginx_apikey_",
		maxAge:    3 * 24 * 3600,
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
	if v.Token() == "" {
		v.New()
	}
	maxAge := time.Duration(s.maxAge) * time.Second
	if len(lifetime) > 0 {
		maxAge = lifetime[0]
	}
	key := s.keyPrefix + v.Token()
	if err := s.client.HSet(ctx, key, v).Err(); err != nil {
		log.Error("Failed to hset", "key", key, "err", err.Error())
		return err
	}
	return s.client.Expire(ctx, key, maxAge).Err()
}
