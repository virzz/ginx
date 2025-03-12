package ginx

import (
	"slices"
	"strings"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/virzz/ginx/rsp"
	"github.com/virzz/vlog"
)

func LogMw(c *gin.Context) {
	c.Next()
	args := []any{
		"remote_ip", c.RemoteIP(),
		"client_ip", c.ClientIP(),
		"referer", c.Request.Referer(),
		"useragent", c.Request.UserAgent(),
		"status", c.Writer.Status(),
	}
	if requestid := requestid.Get(c); requestid != "" {
		args = append(args, "requestid", requestid)
	}
	vlog.Info("AccessLog", args...)
}

func ApikeyMw(name string, apikeys ...string) func(*gin.Context) {
	if name == "" {
		name = "apikey"
	}
	if name == "disable" {
		return func(c *gin.Context) { c.Next() }
	}
	return func(c *gin.Context) {
		apikey := c.GetHeader(name)
		if apikey == "" {
			apikey = c.Query(name)
			if apikey == "" {
				apikey, _ = c.Cookie(name)
			}
		}
		apikey = strings.TrimSpace(apikey)
		if apikey != "" && slices.Contains(apikeys, apikey) {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(401, rsp.M("Error Unauthorized"))
	}
}

func systemAuthMw(token string) func(*gin.Context) {
	return func(c *gin.Context) {
		system := c.GetHeader("system")
		if system == "" {
			system = c.Query("system")
			if system == "" {
				system, _ = c.Cookie("system")
			}
		}
		if system != "" && system == token {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(401, rsp.M("Error Unauthorized"))
	}
}
