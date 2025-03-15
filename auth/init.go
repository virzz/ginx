package auth

import (
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	DefaultKey = "github.com/virzz/ginx/auth"
	TokenKey   = "github.com/virzz/ginx/auth/token"

	paramKey = "token"
)

func IsRole(c *gin.Context, role string) bool {
	return slices.Contains(c.GetStringSlice("roles"), role)
}

func Default[T IDType](c *gin.Context) *Session[T] {
	return c.MustGet(DefaultKey).(*Session[T])
}

func doInit[T IDType](client *redis.Client, data ...Data[T]) gin.HandlerFunc {
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
		if len(token) == 0 {
			token = New()
		}
		var _data Data[T]
		if len(data) > 0 {
			_data = data[0]
			_data.Clear().SetToken(token)
		} else {
			_data = &DefaultData[T]{Token_: token}
		}
		sess := NewSession(c, client, _data)
		c.Set(DefaultKey, sess)
		c.Set(TokenKey, token)

		if !sess.IsNil {
			data := sess.Data()
			roles := data.Roles()
			c.Set("id", data.ID())
			c.Set("account", data.Account())
			c.Set("roles", roles)
			c.Set("is_admin", slices.Contains(roles, "admin"))
			for k, v := range data.Items() {
				c.Set(k, v)
			}
		}
		c.Next()
	}
}

// func Init(cfg *Config) gin.HandlerFunc {
// 	store := redis.NewClient(&redis.Options{Addr: cfg.Addr(), Password: cfg.Pass, DB: cfg.DB})
// 	switch cfg.IDType {
// 	case "int64":
// 		return doInit[int64](store)
// 	case "uint64":
// 		return doInit[uint64](store)
// 	case "string":
// 		return doInit[string](store)
// 	}
// 	panic("idtype unsupported")
// }

func InitWithData[T IDType](cfg *Config, data Data[T]) gin.HandlerFunc {
	store := redis.NewClient(&redis.Options{Addr: cfg.Addr(), Password: cfg.Pass, DB: cfg.DB})
	return doInit(store, data)
}

func Init[T IDType](cfg *Config) gin.HandlerFunc {
	store := redis.NewClient(&redis.Options{Addr: cfg.Addr(), Password: cfg.Pass, DB: cfg.DB})
	return doInit[T](store)
}
