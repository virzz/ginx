package redis

import (
	"context"
	"encoding/base32"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/redis/go-redis/v9"
)

var sessionExpire = 86400 * 30

type RediStore struct {
	Pool          *redis.Client
	Codecs        []securecookie.Codec
	Options       *sessions.Options // default configuration
	DefaultMaxAge int               // default Redis TTL for a MaxAge == 0 session
	maxLength     int
	keyPrefix     string
	serializer    SessionSerializer
}

// Default: 4096,
func (s *RediStore) SetMaxLength(l int)                 { s.maxLength = l }
func (s *RediStore) SetKeyPrefix(p string)              { s.keyPrefix = p }
func (s *RediStore) SetSerializer(ss SessionSerializer) { s.serializer = ss }

func (s *RediStore) SetMaxAge(v int) {
	var c *securecookie.SecureCookie
	var ok bool
	s.Options.MaxAge = v
	for i := range s.Codecs {
		if c, ok = s.Codecs[i].(*securecookie.SecureCookie); ok {
			c.MaxAge(v)
		} else {
			fmt.Printf("Can't change MaxAge on codec %v\n", s.Codecs[i])
		}
	}
}

func NewRediStore(address, password string, keyPairs ...[]byte) (*RediStore, error) {
	s := &RediStore{
		Pool:          redis.NewClient(&redis.Options{Addr: address, Password: password}),
		Options:       &sessions.Options{Path: "/", MaxAge: sessionExpire},
		Codecs:        securecookie.CodecsFromPairs(keyPairs...),
		DefaultMaxAge: 60 * 20, // 20 minutes seems like a reasonable default
		maxLength:     4096,
		keyPrefix:     "session_",
		serializer:    JSONSerializer{},
	}
	return s, s.Pool.Ping(context.Background()).Err()
}

func (s *RediStore) Close() error { return s.Pool.Close() }

func (s *RediStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

func (s *RediStore) New(r *http.Request, name string) (*sessions.Session, error) {
	var (
		err error
		ok  bool
	)
	session := sessions.NewSession(s, name)
	options := *s.Options
	session.Options = &options
	session.IsNew = true
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, s.Codecs...)
		if err == nil {
			ok, err = s.load(r.Context(), session)
			session.IsNew = !(err == nil && ok) // not new if no error and data available
		}
	}
	return session, err
}

func (s *RediStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	// Marked for deletion.
	if session.Options.MaxAge <= 0 {
		if err := s.delete(r.Context(), session); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
	} else {
		if session.ID == "" {
			session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
		}
		if err := s.save(r.Context(), session); err != nil {
			return err
		}
		encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, s.Codecs...)
		if err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	}
	return nil
}

func (s *RediStore) Delete(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	if err := s.delete(r.Context(), session); err != nil {
		return err
	}
	options := *session.Options
	options.MaxAge = -1
	http.SetCookie(w, sessions.NewCookie(session.Name(), "", &options))
	for k := range session.Values {
		delete(session.Values, k)
	}
	return nil
}

func (s *RediStore) delete(ctx context.Context, session *sessions.Session) error {
	return s.Pool.Del(ctx, s.keyPrefix+session.ID).Err()
}

func (s *RediStore) save(ctx context.Context, session *sessions.Session) error {
	b, err := s.serializer.Serialize(session)
	if err != nil {
		return err
	}
	if s.maxLength != 0 && len(b) > s.maxLength {
		return errors.New("SessionStore: the value to store is too big")
	}
	age := session.Options.MaxAge
	if age == 0 {
		age = s.DefaultMaxAge
	}
	return s.Pool.SetEx(ctx, s.keyPrefix+session.ID, b, time.Duration(age)*time.Second).Err()
}

func (s *RediStore) load(ctx context.Context, session *sessions.Session) (bool, error) {
	data, err := s.Pool.Get(ctx, s.keyPrefix+session.ID).Result()
	if data == "" || err == redis.Nil {
		return false, nil // no data was associated with this key
	}
	return true, s.serializer.Deserialize([]byte(data), session)
}
