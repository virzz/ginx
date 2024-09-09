package ginx

import "github.com/gin-gonic/gin"

var (
	mwBefore = []gin.HandlerFunc{}
	mwAfter = []gin.HandlerFunc{}
)
