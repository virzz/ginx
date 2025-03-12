package store

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"time"

	"github.com/redis/go-redis/v9"
)

type Sessioner interface {
	Token() string
	ID() string
	Account() string
	Roles() []string
	Get(string) any
	Set(string, any)
	Clear() error
	Save(...time.Duration) error
}

type Session struct {
	ctx   context.Context
	log   *slog.Logger
	store Store
	data  Data
	IsNil bool
}

var _ Sessioner = (*Session)(nil)

func (s *Session) Token() string           { return s.data.Token() }
func (s *Session) Set(key string, val any) { s.Data().Set(key, val) }
func (s *Session) Get(key string) any      { return s.Data().Get(key) }
func (s *Session) ID() string              { return s.Data().ID() }
func (s *Session) Account() string         { return s.Data().Account() }
func (s *Session) Roles() []string         { return s.Data().Roles() }
func (s *Session) Clear() error            { return s.store.Clear(s.data) }

func (s *Session) Save(lifetime ...time.Duration) error {
	return s.store.Save(s.ctx, s.data, lifetime...)
}

func (s *Session) HasRole(role string) bool {
	return slices.Contains(s.Roles(), role)
}

func (s *Session) Data() Data {
	if s.data.Token() != "" {
		err := s.store.Get(s.ctx, s.data)
		if err != nil {
			s.log.Warn("Failed to get token data", "token", s.data.Token(), "err", err.Error())
			s.IsNil = errors.Is(err, redis.Nil)
		}
	}
	return s.data
}

func NewSession(ctx context.Context, log *slog.Logger, store Store, data Data) *Session {
	return &Session{ctx, log, store, data, false}
}
