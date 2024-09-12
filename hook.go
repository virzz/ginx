package ginx

import "github.com/gin-gonic/gin"

var (
	mwBefore = []gin.HandlerFunc{}
	mwAfter = []gin.HandlerFunc{}
)

func RegisterMwBefore(handler ...gin.HandlerFunc) {
	mwBefore = append(mwBefore, handler...)
}

func RegisterMwAfter(handler ...gin.HandlerFunc) {
	mwAfter = append(mwAfter, handler...)
}
