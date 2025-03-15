package auth

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/virzz/vlog"
)

type RedisStore[T IDType] struct {
	client    redis.UniversalClient
	keyPrefix string
	maxAge    int
}

func NewRedisStore[T IDType](client redis.UniversalClient, maxAge ...int) (*RedisStore[T], error) {
	s := &RedisStore[T]{
		keyPrefix: "ginx:auth:token:",
		maxAge:    7 * 24 * 60 * 60,
		client:    client,
	}
	if len(maxAge) > 0 {
		s.maxAge = maxAge[0]
	}
	return s, client.Ping(context.Background()).Err()
}

func (s *RedisStore[T]) Clear(ctx context.Context, v Data[T]) error {
	return s.client.Del(ctx, s.keyPrefix+v.Token()).Err()
}

func (s *RedisStore[T]) Get(ctx context.Context, v Data[T]) error {
	x := s.client.HGetAll(ctx, s.keyPrefix+v.Token())
	if len(x.Val()) == 0 {
		return redis.Nil
	}
	return x.Scan(v)
}

func (s *RedisStore[T]) Save(ctx context.Context, v Data[T], lifetime ...time.Duration) error {
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
