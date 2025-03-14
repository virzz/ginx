package session

import (
	"github.com/gin-contrib/sessions"
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
		roles, ok := sessions.Default(c).Get("roles").([]any)
		if ok {
			for _, role := range roles {
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
