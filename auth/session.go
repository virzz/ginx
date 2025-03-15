package auth

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/virzz/vlog"
)

type Session[T IDType] struct {
	ctx   context.Context
	data  Data[T]
	IsNil bool

	store     redis.UniversalClient
	keyPrefix string
	maxAge    int
	Token_    string
}

func (s *Session[T]) Token() string            { return s.Token_ }
func (s *Session[T]) ID() T                    { return s.Data().ID() }
func (s *Session[T]) Account() string          { return s.Data().Account() }
func (s *Session[T]) Roles() []string          { return s.Data().Roles() }
func (s *Session[T]) HasRole(role string) bool { return slices.Contains(s.Roles(), role) }

func (s *Session[T]) Data() Data[T] {
	if s.data == nil {
		s.data = &DefaultData[T]{}
		s.data.New()
		return s.data
	}
	if s.Token_ != "" {
		x := s.store.HGetAll(s.ctx, s.keyPrefix+s.Token_)
		if len(x.Val()) == 0 {
			return s.data
		}
		if err := x.Scan(s.data); err != nil {
			vlog.Warn("Failed to get token data", "token", s.Token_, "err", err.Error())
			if !errors.Is(err, redis.Nil) {
				s.Token_ = s.data.Token()
			}
		}
	}
	return s.data
}

func (s *Session[T]) Clear() {
	s.store.Del(s.ctx, s.keyPrefix+s.Token_)
	s.Data().New()
}

func (s *Session[T]) Delete(key string) {
	s.store.HDel(s.ctx, s.keyPrefix+s.Token_, key)
	s.Data().Delete(key)
}

func (s *Session[T]) Get(key string) any            { return s.Data().Get(key) }
func (s *Session[T]) Set(key string, val any)       { s.Data().Set(key, val) }
func (s *Session[T]) SetID(val T)                   { s.Data().SetID(val) }
func (s *Session[T]) SetAccount(val string)         { s.Data().SetAccount(val) }
func (s *Session[T]) SetRoles(roles []string)       { s.Data().SetRoles(roles) }
func (s *Session[T]) SetValues(key string, val any) { s.Data().SetValues(key, val) }
func (s *Session[T]) Save(lifetime ...time.Duration) error {
	if s.Token_ == "" || s.Token_ == "null" {
		s.data.New()
		s.Token_ = s.data.Token()
	}
	maxAge := time.Duration(s.maxAge) * time.Second
	if len(lifetime) > 0 {
		maxAge = lifetime[0]
	}
	key := s.keyPrefix + s.Token_
	if err := s.store.HSet(s.ctx, key, s.data).Err(); err != nil {
		vlog.Error("Failed to hset", "key", key, "err", err.Error())
		return err
	}
	return s.store.Expire(s.ctx, key, maxAge).Err()
}

func NewSession[T IDType](ctx context.Context, store *redis.Client, data Data[T], maxAge ...int) *Session[T] {
	s := &Session[T]{
		ctx:       ctx,
		store:     store,
		data:      data,
		keyPrefix: "ginx:auth:token:",
		maxAge:    7 * 24 * 60 * 60,
		Token_:    data.Token(),
	}
	if len(maxAge) > 0 {
		s.maxAge = maxAge[0]
	}
	return s
}

func (s *Session[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Data())
}

func (s *Session[T]) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, s.Data())
	if err != nil {
		return err
	}
	s.Token_ = s.Data().Token()
	return err
}
