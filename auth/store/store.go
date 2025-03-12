package store

import (
	"context"
	"time"
)

const (
	StoreKeyPrefix = "ginx_session_"
	StoreMaxAge    = 3 * 24 * 3600
)

type (
	Store interface {
		Get(context.Context, Data) error
		Save(context.Context, Data, ...time.Duration) error
		Clear(Data) error

		WithKeyPrefix(string) Store
		WithMaxAge(int) Store
	}

	StoreOption func(Store)
	WithFunc    func(Store) Store
)

func WithKeyPrefix(keyPrefix string) WithFunc {
	return func(r Store) Store { return r.WithKeyPrefix(keyPrefix) }
}

func WithMaxAge(maxAge int) WithFunc {
	return func(r Store) Store { return r.WithMaxAge(maxAge) }
}
