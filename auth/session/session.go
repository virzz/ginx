package session

import (
	"errors"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/context"
)

type Session struct{ sessions.Session }

func (s *Session) Token() string { return s.Session.ID() }

func (s *Session) Get(key string) any            { return s.Session.Get(key) }
func (s *Session) Set(key string, val any)       { s.Session.Set(key, val) }
func (s *Session) SetID(id string)               { s.Set("id", id) }
func (s *Session) SetAccount(account string)     { s.Set("account", account) }
func (s *Session) SetRoles(roles []string)       { s.Set("roles", roles) }
func (s *Session) SetValues(key string, val any) { s.Set(key, val) }
func (s *Session) Save(_ ...time.Duration) error { return s.Session.Save() }

func (s *Session) Delete(key string) { s.Session.Delete(key) }
func (s *Session) Clear()            { s.Session.Clear() }

func (s *Session) getString(key string) string {
	if v := s.Get(key); v != nil {
		if v, ok := v.(string); ok {
			return v
		}
	}
	return ""
}
func (s *Session) ID() string      { return s.getString("id") }
func (s *Session) Account() string { return s.getString("account") }

func (s *Session) HasRole(role string) bool {
	switch roles := s.Get("roles").(type) {
	case []string:
		for _, r := range roles {
			if r == role {
				return true
			}
		}
	case []any:
		for _, r := range roles {
			if v, ok := r.(string); ok && v == role {
				return true
			}
		}
	}
	return false
}
func (s *Session) Roles() []string {
	switch roles := s.Get("roles").(type) {
	case []string:
		return roles
	case []any:
		var rs []string
		for _, r := range roles {
			if v, ok := r.(string); ok {
				rs = append(rs, v)
			}
		}
		return rs
	}
	return nil
}

var ErrNotImplemented = errors.New("not implemented")

func (s *Session) MarshalJSON() ([]byte, error) {
	return nil, ErrNotImplemented
}
func (s *Session) UnmarshalJSON(data []byte) error {
	return ErrNotImplemented
}

func Sessions(name string, store sessions.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := &Session{}
		s.Session = &session{name, c.Request, store, nil, false, c.Writer}
		c.Set(sessions.DefaultKey, s)
		defer context.Clear(c.Request)
		c.Next()
	}
}
