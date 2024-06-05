package apikey

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store interface {
	Get(context.Context, *Data) error
	Save(context.Context, *Data) error
	Clear(*Data) error
}

type RedisStore struct {
	client     redis.UniversalClient // client to connect to redis
	keyPrefix  string                // key prefix with which the session will be stored
	serializer Serializer            // session serializer
	maxAge     int
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
func WithSerializer(serializer Serializer) func(*RedisStore) {
	return func(rs *RedisStore) { rs.serializer = serializer }
}

func NewRedisStore(opts ...Option) (*RedisStore, error) {
	rs := &RedisStore{
		keyPrefix:  "ginx_apikey_",
		serializer: SonicSerializer{},
		maxAge:     3 * 24 * 3600,
	}
	for _, opt := range opts {
		opt(rs)
	}
	if rs.client == nil {
		return nil, errors.New("redisstore: client is nil")
	}
	return rs, rs.client.Ping(context.Background()).Err()
}

func (s *RedisStore) Get(ctx context.Context, v *Data) error  { return s.load(ctx, v) }
func (s *RedisStore) Save(ctx context.Context, v *Data) error { return s.save(ctx, v) }
func (s *RedisStore) Clear(v *Data) error                     { return s.delete(context.Background(), v) }

// func (s *RedisStore) Close() error                            { return s.client.Close() }

func (s *RedisStore) load(ctx context.Context, v *Data) error {
	buf, err := s.client.Get(ctx, s.keyPrefix+v.Token).Bytes()
	if err != nil {
		return err
	}
	return s.serializer.Deserialize(buf, v)
}

func (s *RedisStore) save(ctx context.Context, v *Data) error {
	if v.Token == "" {
		v.New()
	}
	buf, err := s.serializer.Serialize(v)
	if err != nil {
		log.Error("Failed to serialize token data", "err", err.Error())
		return err
	}
	return s.client.Set(ctx, s.keyPrefix+v.Token, buf, time.Duration(s.maxAge)*time.Second).Err()
}

func (s *RedisStore) delete(ctx context.Context, v *Data) error {
	v.Values = nil
	err := s.client.Del(ctx, s.keyPrefix+v.Token).Err()
	if err != nil {
		log.Error("Failed to delete token data", "err", err.Error())
		return err
	}
	return nil
}
