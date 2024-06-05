package ginx

import "github.com/gin-gonic/gin"

type RegisterFunc func(*gin.RouterGroup)

var Routers = []RegisterFunc{}

func Register(registers ...RegisterFunc) {
	Routers = append(Routers, registers...)
}
