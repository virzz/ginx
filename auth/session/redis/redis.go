package redis

import (
	"errors"

	"github.com/gin-contrib/sessions"
)

type Store interface {
	sessions.Store
}

type store struct {
	*RediStore
}

func NewStore(address, password string, keyPairs ...[]byte) (Store, error) {
	s, err := NewRediStore(address, password, keyPairs...)
	if err != nil {
		return nil, err
	}
	return &store{s}, nil
}

func GetRedisStore(s Store) (err error, rediStore *RediStore) {
	realStore, ok := s.(*store)
	if !ok {
		err = errors.New("unable to get the redis store: Store isn't *store")
		return
	}
	rediStore = realStore.RediStore
	return
}

func SetKeyPrefix(s Store, prefix string) error {
	err, rediStore := GetRedisStore(s)
	if err != nil {
		return err
	}
	rediStore.SetKeyPrefix(prefix)
	return nil
}

func (c *store) Options(options sessions.Options) {
	c.RediStore.Options = options.ToGorillaOptions()
}
