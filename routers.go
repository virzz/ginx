package ginx

import "github.com/gin-gonic/gin"

type (
	Controller interface {
		Authed(*gin.RouterGroup)
		UnAuth(*gin.RouterGroup)
	}

	RegisterFunc func(*gin.RouterGroup)
)

var Routers = []RegisterFunc{}

func Register(vs ...RegisterFunc) { Routers = append(Routers, vs...) }
