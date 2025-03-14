package token

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

func Default(c *gin.Context) *Session {
	return c.MustGet(DefaultKey).(*Session)
}

func IsRole(c *gin.Context, role string) bool {
	return slices.Contains(c.GetStringSlice("roles"), role)
}

func Init(s *redis.Client, data ...Data) gin.HandlerFunc {
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
		sess := NewSession(c, s, &DefaultData{Token_: token})
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
