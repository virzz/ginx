package ginx

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/virzz/ginx/rsp"
	"github.com/virzz/vlog"
)

func LogMw(c *gin.Context) {
	vlog.Info("AccessLog",
		"remote_ip", c.RemoteIP(),
		"client_ip", c.ClientIP(),
		"referer", c.Request.Referer(),
		"useragent", c.Request.UserAgent(),
	)
	c.Next()
}

func AuthMw(apikeys ...string) func(*gin.Context) {
	return func(c *gin.Context) {
		apikey := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if apikey == "" {
			apikey = c.Query("apikey")
			if apikey == "" {
				apikey, _ = c.Cookie("apikey")
			}
		}
		if apikey != "" {
			for _, key := range apikeys {
				if apikey == key {
					c.Next()
					return
				}
			}
		}
		c.AbortWithStatusJSON(401, rsp.M("Error Unauthorized"))
	}
}
