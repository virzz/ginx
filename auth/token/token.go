package token

import (
	"log/slog"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/virzz/vlog"

	"github.com/virzz/ginx/auth/store"
	"github.com/virzz/ginx/code"
	"github.com/virzz/ginx/rsp"
)

const (
	DefaultKey = "github.com/virzz/ginx/auth"
	TokenKey   = "github.com/virzz/ginx/auth/token"

	paramKey = "token"
)

var log *slog.Logger

func Init(s store.Store, data ...store.Data) gin.HandlerFunc {
	log = vlog.Log.WithGroup("token")
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
		var _data store.Data
		if len(data) > 0 {
			_data = data[0]
		} else {
			_data = &store.DefaultData{}
		}
		_data.SetToken(token)
		c.Set(DefaultKey, store.NewSession(c.Copy(), log, s, _data))
		c.Set(TokenKey, token)
		c.Next()
	}
}

func Default(c *gin.Context) *store.Session {
	return c.MustGet(DefaultKey).(*store.Session)
}

func AuthedMW(c *gin.Context) {
	sess := Default(c)
	if sess.Token() == "" {
		c.AbortWithStatusJSON(200, rsp.C(code.TokenInvalid))
		return
	}
	data := sess.Data()
	if sess.IsNil {
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
