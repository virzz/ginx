package apikey

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/virzz/ginx/code"
	"github.com/virzz/ginx/rsp"
	"github.com/virzz/vlog"
)

const (
	DefaultKey = "github.com/virzz/gin-apikey"
	TokenKey   = "github.com/virzz/gin-apikey/token"

	paramKey = "token"
)

var log *slog.Logger

type APIKey interface {
	Token() string
	ID() string
	Account() string
	Roles() []string
	Get(string) any
	Set(string, any)
	Clear() error
	Save(...time.Duration) error
}

type session struct {
	ctx   context.Context
	store Store
	data  Data
	isNil bool
}

var _ APIKey = (*session)(nil)

func (s *session) Token() string           { return s.data.Token() }
func (s *session) Set(key string, val any) { s.Data().Set(key, val) }
func (s *session) Get(key string) any      { return s.Data().Get(key) }
func (s *session) ID() string              { return s.Data().ID() }
func (s *session) Account() string         { return s.Data().Account() }
func (s *session) Roles() []string         { return s.Data().Roles() }
func (s *session) Clear() error            { return s.store.Clear(s.data) }

func (s *session) Save(lifetime ...time.Duration) error {
	return s.store.Save(s.ctx, s.data, lifetime...)
}

func (s *session) HasRole(role string) bool {
	return slices.Contains(s.Roles(), role)
}

func (s *session) Data() Data {
	if s.data.Token() != "" {
		err := s.store.Get(s.ctx, s.data)
		if err != nil {
			log.Warn("Failed to get token data", "token", s.data.Token(), "err", err.Error())
			s.isNil = errors.Is(err, redis.Nil)
		}
	}
	return s.data
}

func Init(store Store, data ...Data) gin.HandlerFunc {
	log = vlog.Log.WithGroup("apikey")
	return func(c *gin.Context) {
		token := c.Query(paramKey)
		if len(token) == 0 {
			token = c.PostForm(paramKey)
			if len(token) == 0 {
				token, _ = c.Cookie(paramKey)
				if len(token) == 0 {
					token = strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
				}
			}
		}
		var _data Data
		if len(data) > 0 {
			_data = data[0]
		} else {
			_data = &DefaultData{}
		}
		_data.SetToken(token)
		c.Set(DefaultKey, &session{c.Copy(), store, _data, false})
		c.Set(TokenKey, token)
		c.Next()
	}
}

func Default(c *gin.Context) *session { return c.MustGet(DefaultKey).(*session) }

func AuthedMW(c *gin.Context) {
	sess := Default(c)
	if sess.Token() == "" {
		c.AbortWithStatusJSON(200, rsp.C(code.TokenInvalid))
		return
	}
	data := sess.Data()
	fmt.Println(data, sess.isNil)
	if sess.isNil {
		c.AbortWithStatusJSON(200, rsp.C(code.TokenExpired))
		return
	}
	if data.ID() == "" {
		c.AbortWithStatusJSON(200, rsp.C(code.Forbidden))
		return
	}
	c.Set("id", data.ID())
	c.Set("account", data.Account())
	c.Set("roles", data.Roles())
	c.Set("is_admin", slices.Contains(data.Roles(), "admin"))
	c.Next()
}

func IsRole(c *gin.Context, role string) bool {
	if v, ok := c.Get("roles"); ok {
		if vv, ok := v.([]string); ok {
			return slices.Contains(vv, role)
		}
	}
	return false
}

func AuthRoleMW(roles ...string) gin.HandlerFunc {
	roleMap := make(map[string]struct{})
	for _, role := range roles {
		roleMap[role] = struct{}{}
	}
	return func(c *gin.Context) {
		for _, r := range Default(c).Roles() {
			if _, ok := roleMap[r]; ok {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(200, rsp.C(code.Forbidden))
	}
}
