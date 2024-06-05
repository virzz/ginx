package apikey

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"strings"

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
	Values() map[string]any
	Get(key string) any
	Set(key string, val any)
	Delete(key string)
	Clear() error
	Save() error
}

type session struct {
	ctx   context.Context
	store Store
	data  *Data
	isNil bool
}

var _ APIKey = (*session)(nil)

func (s *session) Get(key string) any      { return s.Data().Values[key] }
func (s *session) Set(key string, val any) { s.Data().Values[key] = val }
func (s *session) Delete(key string)       { delete(s.Data().Values, key) }
func (s *session) Token() string           { return s.Data().Token }
func (s *session) Values() map[string]any  { return s.Data().Values }
func (s *session) Save() error             { return s.store.Save(s.ctx, s.data) }
func (s *session) Clear() error            { return s.store.Clear(s.data) }
func (s *session) Roles() []string         { return s.Data().Values["roles"].([]string) }
func (s *session) HasRole(role string) bool {
	roles, ok := s.Get("roles").([]string)
	return ok && slices.Contains(roles, role)
}
func (s *session) Data() *Data {
	if len(s.data.Values) > 0 {
		return s.data
	}
	if s.data.Token != "" {
		err := s.store.Get(s.ctx, s.data)
		if err != nil {
			log.Warn("Failed to get token data", "token", s.data.Token, "err", err.Error())
			s.isNil = errors.Is(err, redis.Nil)
		}
	}
	if s.data.Values == nil {
		s.data.Values = make(map[string]any)
	}
	return s.data
}

func Init(store Store) gin.HandlerFunc {
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
		c.Set(TokenKey, token)
		s := &session{c.Copy(), store, &Data{Token: token}, false}
		c.Set(DefaultKey, s)
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
	if sess.Data(); sess.isNil {
		c.AbortWithStatusJSON(200, rsp.C(code.TokenExpired))
		return
	}
	if v, ok := sess.Get("id").(string); !ok || v == "" {
		c.AbortWithStatusJSON(200, rsp.C(code.Forbidden))
		return
	}
	for k, v := range sess.Values() {
		c.Set(k, v)
	}
	c.Next()
}

func AuthRoleMW(roles ...string) gin.HandlerFunc {
	roleMap := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		roleMap[role] = struct{}{}
	}
	return func(c *gin.Context) {
		if vs, ok := Default(c).Get("roles").([]any); ok {
			for _, role := range vs {
				if r, ok := role.(string); ok {
					if _, ok := roleMap[r]; ok {
						c.Next()
						return
					}
				}
			}
		}
		c.AbortWithStatusJSON(200, rsp.C(code.Forbidden))
	}
}
