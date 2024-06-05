package ginx

import "github.com/gin-gonic/gin"

type Controller interface {
	Authed(g *gin.RouterGroup)
	UnAuth(g *gin.RouterGroup)
}
