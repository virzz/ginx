package token

import (
	"github.com/gin-gonic/gin"
	"github.com/virzz/ginx/code"
	"github.com/virzz/ginx/rsp"
)

func RoleMW(roles ...string) gin.HandlerFunc {
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

func AuthMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("id") != "" && c.GetString("account") != "" {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(200, rsp.C(code.Unauthorized))
	}
}
