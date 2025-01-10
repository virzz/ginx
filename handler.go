package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/virzz/ginx/apikey"
	"github.com/virzz/ginx/code"
	"github.com/virzz/ginx/rsp"
)

func LogoutHandler(c *gin.Context) {
	err := apikey.Default(c).Clear()
	if err != nil {
		c.AbortWithStatusJSON(200, rsp.C(code.TokenDestory))
		return
	}
	c.JSON(200, rsp.OK())
}

func CodesHandler(c *gin.Context) {
	c.JSON(200, rsp.S(code.Codes))
}

func HealthCheckHandler(c *gin.Context) {
	c.Status(200)
}
