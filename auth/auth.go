package auth

import (
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/virzz/ginx/auth/session"
	redisStore "github.com/virzz/ginx/auth/session/redis"
	"github.com/virzz/ginx/auth/token"
	"github.com/virzz/vlog"
)

type Session interface {
	Token() string

	ID() string
	Get(string) any
	Account() string
	Roles() []string
	HasRole(string) bool

	Delete(string)
	Clear()

	Set(string, any)
	SetAccount(string)
	SetID(string)
	SetRoles([]string)
	SetValues(string, any)
	Save(...time.Duration) error
}

var _ Session = (*token.Session)(nil)
var _ Session = (*session.Session)(nil)

func Default(c *gin.Context) Session {
	v, ok := c.Get(token.DefaultKey)
	if ok {
		return v.(Session)
	}
	v, ok = c.Get(sessions.DefaultKey)
	if ok {
		vlog.Infof("%+v", v)
		vlog.Info("sessions")
		return v.(Session)
	}
	panic("no session found")
}

func Init(c *Config) gin.HandlerFunc {
	switch c.Type {
	case AuthTypeToken:
		client := redis.NewClient(&redis.Options{Addr: c.Addr(), Password: c.Pass, DB: c.DB})
		return token.Init(client)
	case AuthTypeSession:
		store, err := redisStore.NewStore(c.Addr(), c.Pass, []byte(c.Secret))
		if err != nil {
			panic(err)
		}
		return session.Sessions("session", store)
	case AuthTypeCookie:
		store := cookie.NewStore([]byte(c.Secret))
		return session.Sessions("session", store)
	}
	return nil
}
